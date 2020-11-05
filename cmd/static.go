/*
Copyright © 2020 Ulrich Wisser <ulrich@wisser.se>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"log"
	"strconv"
	"time"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/DNS-OARC/ripeatlas/measurement/dns"
	client "github.com/influxdata/influxdb/client/v2"
	mdns "github.com/miekg/dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var statsStaticRcode [32]int
var statsStaticNSEC int
var statsStaticNSEC3 int
var statsStaticNONSEC int

// staticCmd represents the static command
var staticCmd = &cobra.Command{
	Use:   "static",
	Short: "count rcodes and nsec types for statis non existent domain",
	Long:  `count rcodes and nsec types for statis non existent domain`,
	Run:   runStatic,
}

func init() {
	rootCmd.AddCommand(staticCmd)

	// Use flags for viper values
	viper.BindPFlags(invalidCmd.Flags())
}

func runStatic(cmd *cobra.Command, args []string) {
	var measurements = make([]int, 0)

	// check config
	checkInvalidConf()

	// check arguments
	if len(args) == 0 {
		log.Fatal("At least one measurement id must be given")
	}

	// convert arguments
	for _, m := range args {
		v, err := strconv.Atoi(m)
		if err != nil {
			log.Fatal("Could not convert to int: ", m)
		}
		measurements = append(measurements, v)
	}

	ch := subscribe(measurements)

	go rcvStatic(ch)
	go rcvStatic(ch)

	go saveStaticStats()

	select {}
}

func checkStaticConf() {
	if len(viper.GetString("influxserver")) == 0 {
		log.Println("Influx server must be given")
	}
	if len(viper.GetString("influxdb")) == 0 {
		log.Println("Influx database must be given")
	}
}

func rcvStatic(ch <-chan *measurement.Result) {
	for {
		// rcv measurement
		msm := <-ch

		// if parsing fails
		if msm.ParseError != nil {
			log.Println(msm.ParseError.Error())
			continue
		}

		// we handle only dns results
		if msm.Type() != "dns" {
			log.Printf("Wrong result type msmid %d type %s", msm.MsmId(), msm.Type())
			return
		}

		// debug output of received result
		if verbose > 2 {
			log.Printf("%d %s", msm.MsmId(), msm.Type())
		}

		// handle single result
		if msm.DnsResult() != nil {
			handleStatic(msm.DnsResult())
		}
		for _, s := range msm.DnsResultsets() {
			if s.Result() != nil {
				handleStatic(s.Result())
			}
		}
	}

}

func handleStatic(result *dns.Result) {

	msg, err := result.UnpackAbuf()
	if err != nil {
		log.Println("Could not unpack Abuf ", err)
		return
	}
	statsStaticRcode[msg.Rcode]++
	switch nsec(msg.Ns) {
	case NSEC:
		statsStaticNSEC++
	case NSEC3:
		statsStaticNSEC3++
	case NONSEC:
		statsStaticNONSEC++
	}
}

func saveStaticStats() {
	// influxdb client config
	conf := client.HTTPConfig{
		Addr: viper.GetString("influxserver"),
	}
	if len(viper.GetString("influxuser")) > 0 {
		conf.Username = viper.GetString("influxuser")
		conf.Password = viper.GetString("influxpasswd")
	}

	// get access to Influx
	influx, err := client.NewHTTPClient(conf)
	if err != nil {
		log.Fatalf("Error creating InfluxDB Client: %s", err.Error())
	}

	ticker := time.NewTicker(15 * time.Second)
	for {
		now := <-ticker.C

		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  viper.GetString("influxdb"),
			Precision: "s",
		})
		if err != nil {
			log.Fatalf("Could not create new batch points: %s", err.Error())
		}

		// tags
		tags := map[string]string{}

		// values
		fields := map[string]interface{}{}
		for _, rcode := range RCODES {
			fields[mdns.RcodeToString[rcode]] = statsStaticRcode[rcode]
		}

		// create new point
		pt, err := client.NewPoint("staticRcodes", tags, fields, now)
		if err != nil {
			log.Fatal("Could not create new point. ", err)
		}
		if verbose > 2 {
			log.Printf("Tags:   %v\n", tags)
			log.Printf("Fields: %v\n", fields)
		}

		// add point to list
		bp.AddPoint(pt)

		// values
		fields = map[string]interface{}{
			"nsec":   statsStaticNSEC,
			"nsec3":  statsStaticNSEC3,
			"nonsec": statsStaticNONSEC,
		}

		// create new point
		pt, err = client.NewPoint("staticNsec", tags, fields, now)
		if err != nil {
			log.Fatal("Could not create new point. ", err)
		}
		if verbose > 2 {
			log.Printf("Tags:   %v\n", tags)
			log.Printf("Fields: %v\n", fields)
		}

		// add point to list
		bp.AddPoint(pt)

		// write to database
		if verbose > 1 {
			log.Println("Writing to influx")
		}
		err = influx.Write(bp)
		if err != nil {
			log.Printf("Error writing to InfluxDB: %s", err.Error())
		}
	}
}
