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

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"

	mdns "github.com/miekg/dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var statsRandomRcode [32]int
var statsRandomNSEC int
var statsRandomNSEC3 int
var statsRandomNONSEC int

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "count rcodes and nsec for random domain",
	Long:  `count rcodes and nsec for random domain`,
	Run:   runRandom,
}

func init() {
	rootCmd.AddCommand(randomCmd)

	// Use flags for viper values
	viper.BindPFlags(randomCmd.Flags())
}

func runRandom(cmd *cobra.Command, args []string) {
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

	go rcvRandom(ch)
	go rcvRandom(ch)

	go saveRandomStats()

	select {}
}

func checkRandomConf() {
	if len(viper.GetString("influxserver")) == 0 {
		log.Println("Influx server must be given")
	}
	if len(viper.GetString("influxdb")) == 0 {
		log.Println("Influx database must be given")
	}
}

func rcvRandom(ch <-chan *measurement.Result) {
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
			handleRandom(msm.DnsResult())
		}
		for _, s := range msm.DnsResultsets() {
			if s.Result() != nil {
				handleRandom(s.Result())
			}
		}
	}

}

func handleRandom(result *dns.Result) {

	msg, err := result.UnpackAbuf()
	if err != nil {
		log.Println("Could not unpack Abuf ", err)
		return
	}
	statsRandomRcode[msg.Rcode]++
	switch nsec(msg.Ns) {
	case NSEC:
		statsRandomNSEC++
	case NSEC3:
		statsRandomNSEC3++
	case NONSEC:
		statsRandomNONSEC++
	}
}

func saveRandomStats() {
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

	ticker := time.NewTicker(UPDATEINTERVAL)
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
			fields[mdns.RcodeToString[rcode]] = statsRandomRcode[rcode]
		}

		// create new point
		pt, err := client.NewPoint("randomRcodes", tags, fields, now)
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
			"nsec":   statsRandomNSEC,
			"nsec3":  statsRandomNSEC3,
			"nonsec": statsRandomNONSEC,
		}

		// create new point
		pt, err = client.NewPoint("randomNsec", tags, fields, now)
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
