package cmd

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

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var start time.Time

// measureCmd represents the measure command
var measureCmd = &cobra.Command{
	Use:   "measure",
	Short: "Start a measurement",
	Long:  `Starts a Ripe Atlas measurement of NSEC answers.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Use flags for viper values
		viper.BindPFlags(cmd.Flags())
		checkConf()

		// compute start and stop time
		d, _ := cmd.Flags().GetDuration("duration")
		stop := start.Add(d)

		req1 := &MeasurementRequest{
			Definitions: makeDefinitions(),
			BillTo:      viper.GetString("RIPEACCOUNT"),
			IsOneoff:    false,
			Probes:      probes,
			StartTime:   int(start.Unix()),
			StopTime:    int(stop.Unix()),
		}

		resp1 := createMeasurement(req1)
		log.Printf("INVALID/STATIC/RANDOM %v", resp1)

		req2 := &MeasurementRequest{
			Definitions: makeAuth4Definitions(),
			BillTo:      viper.GetString("RIPEACCOUNT"),
			IsOneoff:    false,
			Probes:      probesV4,
			StartTime:   int(start.Unix()),
			StopTime:    int(stop.Unix()),
		}

		resp2 := createMeasurement(req2)
		log.Printf("AUTHORITATIVE v4 %v", resp2)

		req3 := &MeasurementRequest{
			Definitions: makeAuth6Definitions(),
			BillTo:      viper.GetString("RIPEACCOUNT"),
			IsOneoff:    false,
			Probes:      probesV6,
			StartTime:   int(start.Unix()),
			StopTime:    int(stop.Unix()),
		}

		resp3 := createMeasurement(req3)
		log.Printf("AUTHORITATIVE  v6 %v", resp3)
	},
}

func init() {
	rootCmd.AddCommand(measureCmd)
	measureCmd.Flags().StringP("invalid", "i", "", "domain name with invalid dnssec chain")
	measureCmd.Flags().StringP("static", "s", "", "domain name that is not in the zone")
	measureCmd.Flags().StringP("random", "r", "", "base name for random, static tests ")
	measureCmd.Flags().StringP("begin", "b", "", "time and date for the measurement to start (empty=now)")
	measureCmd.Flags().StringSliceP("authoritative", "a", []string{}, "name of authoritative name servers")
	measureCmd.Flags().DurationP("duration", "d", 4*time.Hour, "how long the measurement should be run")

	// Use flags for viper values
	viper.BindPFlags(measureCmd.Flags())
}

func checkConf() {
	// get command line arguments
	i := viper.GetString("invalid")
	n := viper.GetString("static")
	r := viper.GetString("random")
	s := viper.GetString("begin")
	d := viper.GetDuration("duration")

	// get conf file arguments
	a := viper.GetString("RIPEACCOUNT")
	k := viper.GetString("APIKEY")

	// check arguments
	if len(i) == 0 {
		log.Fatal("invalid must be given")
	}
	if len(n) == 0 {
		log.Fatal("static must be given")
	}
	if len(r) == 0 {
		log.Fatal("random must be given")
	}

	if len(s) == 0 {
		d, _ := time.ParseDuration("1m")
		start = time.Now().Add(d)
	} else {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			log.Fatal("Could not parse start time. ", err)
		}
		start = t
	}

	if len(a) == 0 {
		log.Fatal("RIPEACCOUNT must be given in config file")
	}

	if len(k) == 0 {
		log.Fatal("APIKEY must be given in config file")
	}

	// debug output
	if verbose > 0 {
		log.Println("Invalid:       ", i)
		log.Println("Static:      ", n)
		log.Println("Random:        ", r)
		log.Println("Start:         ", start)
		log.Println("Duration:      ", d.String())
		log.Println("Authoritative: ", viper.GetStringSlice("authoritative"))
		log.Println("Ripe account:  ", viper.GetString("RIPEACCOUNT"))
		log.Println("APIKEY:        ", viper.GetString("APIKEY"))
	}
}

func makeDefinitions() []Definition {
	defs := make([]Definition, 0)

	invalid := viper.GetString("invalid")
	static := viper.GetString("static")
	random := viper.GetString("random")

	// an invalid signed domain
	// if we get data, we know that DNSSEC is not validated
	def1 := definition1
	def1.Description = "Invalid signed domain"
	def1.QueryArgument = invalid
	defs = append(defs, def1)

	// qyuery same name all the time
	// should see when caches expire
	def2 := definition1
	def2.Description = "Static nxdomain"
	def2.QueryArgument = static
	defs = append(defs, def2)

	// query for random domain
	// result should be static (not servfail)
	def3 := definition1
	def3.UseMacros = true
	def3.Description = "Random nxdomain"
	def3.QueryArgument = "$r-$p-$t-" + random
	defs = append(defs, def3)

	// done
	return defs
}
func makeAuth4Definitions() []Definition {
	defs := make([]Definition, 0)

	static := viper.GetString("static")

	// definitions for authoritative servers
	for _, ns := range viper.GetStringSlice("authoritative") {
		def4 := definition2
		def4.AF = 4
		def4.QueryArgument = static
		def4.Target = ns
		defs = append(defs, def4)
	}

	// done
	return defs
}
func makeAuth6Definitions() []Definition {
	defs := make([]Definition, 0)

	static := viper.GetString("static")

	// definitions for authoritative servers
	for _, ns := range viper.GetStringSlice("authoritative") {
		def6 := definition2
		def6.AF = 6
		def6.QueryArgument = static
		def6.Target = ns
		defs = append(defs, def6)
	}

	// done
	return defs
}

// createMeasurement creates a measurement for all types
func createMeasurement(d *MeasurementRequest) *MeasurementResponse {
	apiurl, _ := url.Parse(fmt.Sprintf("https://atlas.ripe.net/api/v2/measurements/?key=%s", viper.GetString("APIKEY")))

	// prepare request
	req, err := http.NewRequest(http.MethodPost, apiurl.String(), nil)
	if err != nil {
		log.Fatal("error setting up request. ", err)
	}
	req.Header.Set("User-Agent", "nsecmonitor/0.0")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// prepare request body
	body, err := json.Marshal(d)
	if err != nil {
		log.Fatalf("error marshalling: %v", d)
	}
	//log.Println(string(body))
	buf := bytes.NewReader(body)
	req.Body = ioutil.NopCloser(buf)
	req.ContentLength = int64(buf.Len())
	if verbose > 1 {
		//log.Println(formatRequest(req))
	}
	// setup http client
	client := &http.Client{Timeout: 20 * time.Second}

	// make api call
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("HTTP call error: %v", err)
	}

	// Everything is fine
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Fatal("API call returned status ", resp.Status)
	}

	// prepare return value
	m := &MeasurementResponse{}

	// If there is a body, get it
	if resp.Body != nil {
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("could not read body: ", err)
		}
		defer resp.Body.Close()
	}

	err = json.Unmarshal(body, m)
	if err != nil {
		var apierror APIError
		err = json.Unmarshal(body, &apierror)
		if err != nil {
			log.Fatal("Could not unmarshal body")
		}
		log.Fatalf("RIPE API call fails with: %v", apierror)

	}

	return m
}
