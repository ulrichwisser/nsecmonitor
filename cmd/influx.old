package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/influxdata/influxdb1-client"

	// this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
)

func write2influx(config *Configuration, responses []Response) {
	// get access to Influx
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.InfluxServer,
		Username: config.InfluxUser,
		Password: config.InfluxPasswd,
	})
	if err != nil {
		log.Fatalf("Error creating InfluxDB Client: %s", err.Error())
	}
	defer c.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  config.InfluxDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatalf("Could not create new batch points: %s", err.Error())
	}

	now := time.Now()
	lineSeen := make(map[string]bool)
	for _, response := range responses {
		// filter wrong direction
		if response.JourneyDirection != config.SiteDirection {
			continue
		}
		// filter already seen lines
		if lineSeen[response.LineNumber] {
			continue
		}

		// parse timestamps
		dtTimeTable, _ := time.Parse("2006-01-02T15:04:05", response.TimeTabledDateTime)
		dtExpected, _ := time.Parse("2006-01-02T15:04:05", response.ExpectedDateTime)

		// tags
		tags := map[string]string{
			"linenumber":    response.LineNumber,
			"destination":   response.Destination,
			"siteid":        config.SiteID,
			"sitedirection": fmt.Sprintf("%d", response.JourneyDirection),
		}

		// values
		fields := map[string]interface{}{
			"TimeTableDateTime":  response.TimeTabledDateTime,
			"ExpectedDateTime":   response.ExpectedDateTime,
			"DisplayTime":        response.DisplayTime,
			"TimeTableTimestamp": dtTimeTable.Unix(),
			"ExpectedTimestamp":  dtExpected.Unix(),
		}

		// create new point
		pt, err := client.NewPoint("SLlatest", tags, fields, now)
		if err != nil {
			log.Fatal(err)
		}
		if config.Verbose {
			log.Printf("Tags:   %v\n", tags)
			log.Printf("Fields: %v\n", fields)
			log.Printf("Point:  %v\n", pt)
		}

		// prevent data for same line
		lineSeen[response.LineNumber] = true

		// add point to list
		bp.AddPoint(pt)
	}

	// Write the batch
	if config.Dryrun {
		log.Println("DRYRUN - NO DATA WILL BE WRITTEN TO INFLUXDB")
	} else {
		err = c.Write(bp)
		if err != nil {
			log.Fatalf("Error writing to InfluxDB: %s", err.Error())
		}
	}
}
