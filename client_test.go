package influxdbhelper

import (
	"time"
)

func ExampleClient_WritePoint() {
	c, _ := NewClient("http://localhost:8086", "", "", "ns")

	type EnvSample struct {
		Time        time.Time `influx:"time"`
		Location    string    `influx:"location,tag"`
		Temperature float64   `influx:"temperature"`
		Humidity    float64   `influx:"humidity"`
		Id          string    `influx:"-"`
	}

	s := EnvSample{
		Time:        time.Now(),
		Location:    "Rm 243",
		Temperature: 70.0,
		Humidity:    60.0,
		Id:          "12432as32",
	}

	c.UseDB("myDb").UseMeasurement("test").WritePoint(s)
}

func ExampleClient_Query() {
	c, _ := NewClient("http://localhost:8086", "", "", "ns")

	type EnvSample struct {
		Time        time.Time `influx:"time"`
		Location    string    `influx:"location,tag"`
		Temperature float64   `influx:"temperature"`
		Humidity    float64   `influx:"humidity"`
		Id          string    `influx:"-"`
	}

	samplesRead := []EnvSample{}

	q := `SELECT * FROM test ORDER BY time DESC LIMIT 10`

	c.UseDB("myDb").DecodeQuery(q, &samplesRead)

	// samplesRead is now populated with data from InfluxDb
}
