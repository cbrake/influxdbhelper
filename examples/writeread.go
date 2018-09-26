package main

import (
	"log"
	"time"

	"github.com/cbrake/influxdbhelper"
	client "github.com/influxdata/influxdb/client/v2"
)

const (
	// Static connection configuration
	influxURL = "http://localhost:8086"
	db        = "dbhelper"
)

var c *influxdbhelper.Client

// Init initializes the database connection
func Init() (err error) {
	c, err = influxdbhelper.NewClient(influxURL, "", "", "ns")
	if err != nil {
		return
	}
	// Create MM database if it doesn't already exist
	q := client.NewQuery("CREATE DATABASE "+db, "", "")
	res, err := c.InfluxClient().Query(q)
	if err != nil {
		return err
	}
	if res.Error() != nil {
		return res.Error()
	}
	log.Println("dbhelper db initialized")
	return nil
}

type envSample struct {
	Time        time.Time `influx:"time"`
	Location    string    `influx:"location,tag"`
	Temperature float64   `influx:"temperature"`
	Humidity    float64   `influx:"humidity"`
	ID          string    `influx:"-"`
}

// we populate a few more fields when reading back
// date to verify unused fields are handled correctly
type envSampleRead struct {
	Time        time.Time `influx:"time"`
	Location    string    `influx:"location,tag"`
	City        string    `influx:"city,tag"`
	Temperature float64   `influx:"temperature"`
	Humidity    float64   `influx:"humidity"`
	Cycles      float64   `influx:"cycles"`
	ID          string    `influx:"-"`
}

func generateSampleData() []envSample {
	ret := make([]envSample, 10)

	for i := range ret {
		ret[i] = envSample{
			Time:        time.Now(),
			Location:    "Rm 243",
			Temperature: 70 + float64(i),
			Humidity:    60 - float64(i),
			ID:          "12432as32",
		}
	}

	return ret
}

func main() {
	err := Init()
	if err != nil {
		log.Fatal("Failed to initialize db")
	}

	// write sample data to database
	samples := generateSampleData()
	for _, p := range samples {
		err := c.WritePoint(db, "test", p)
		if err != nil {
			log.Fatal("Error writing point: ", err)
		}
	}

	// query data from db
	samplesRead := []envSampleRead{}

	q := `SELECT * FROM test ORDER BY time DESC LIMIT 10`
	err = c.Query(db, q, &samplesRead)
	if err != nil {
		log.Fatal("Query error: ", err)
	}
	log.Printf("Samples read: %+v\n", samplesRead)
}
