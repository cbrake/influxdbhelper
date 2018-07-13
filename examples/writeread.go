package main

import (
	"log"
	"time"

	"github.com/cbrake/influxdbhelper"
	client "github.com/influxdata/influxdb/client/v2"
)

const (
	// Static connection configuration
	influxUrl = "http://localhost:8086"
	db        = "dbhelper"
)

var c *influxdbhelper.Client

func Init() (err error) {
	c, err = influxdbhelper.NewClient(influxUrl, "", "", "ns")
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

type EnvSample struct {
	Time        time.Time `influx:"time"`
	Location    string    `influx:"location,tag"`
	Temperature float64   `influx:"temperature"`
	Humidity    float64   `influx:"humidity"`
	Id          string    `influx:"-"`
}

func generateSampleData() []EnvSample {
	ret := make([]EnvSample, 10)

	for i, _ := range ret {
		ret[i] = EnvSample{
			Time:        time.Now(),
			Location:    "Rm 243",
			Temperature: 70 + float64(i),
			Humidity:    60 - float64(i),
			Id:          "12432as32",
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
	samplesRead := []EnvSample{}

	q := `SELECT * FROM test ORDER BY time DESC LIMIT 10`
	err = c.Query(db, q, &samplesRead)
	if err != nil {
		log.Fatal("Query error: ", err)
	}
	log.Printf("Samples read: %+v\n", samplesRead)
}
