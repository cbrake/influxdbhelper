# InfluxDb Helper Library

> easily write and query InfluxDb from Go programs

This library allows you to encode/decode InfluxDb data to/from
Go structs -- similiar to JSON and MongoDb using Go struct field tags.

See [GoDoc](https://godoc.org/github.com/cbrake/influxdbhelper) for more documentation.

## Install

```
go get github.com/influxdata/influxdb1-client
go get github.com/cbrake/influxdbhelper
```

Note, this library currently does not work with the 1.7.x version of the InfluxDB client. It has
been tested with 1.6.5.

## Example

```go
package main

import (
	"log"
	"time"

	"github.com/cbrake/influxdbhelper"
	client "github.com/influxdata/influxdb1-client/v2"
)

const (
	// Static connection configuration
	influxURL = "http://localhost:8086"
	db        = "dbhelper"
)

var c influxdbhelper.Client

// Init initializes the database connection
func Init() (err error) {
	c, err = influxdbhelper.NewClient(influxURL, "", "", "ns")
	if err != nil {
		return
	}
	// Create test database if it doesn't already exist
	q := client.NewQuery("CREATE DATABASE "+db, "", "")
	res, err := c.Query(q)
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
	InfluxMeasurement influxdbhelper.Measurement
	Time              time.Time `influx:"time"`
	Location          string    `influx:"location,tag"`
	Temperature       float64   `influx:"temperature"`
	Humidity          float64   `influx:"humidity"`
	ID                string    `influx:"-"`
}

// we populate a few more fields when reading back
// date to verify unused fields are handled correctly
type envSampleRead struct {
	InfluxMeasurement influxdbhelper.Measurement
	Time              time.Time `influx:"time"`
	Location          string    `influx:"location,tag"`
	City              string    `influx:"city,tag,field"`
	Temperature       float64   `influx:"temperature"`
	Humidity          float64   `influx:"humidity"`
	Cycles            float64   `influx:"cycles"`
	ID                string    `influx:"-"`
}

func generateSampleData() []envSample {
	ret := make([]envSample, 10)

	for i := range ret {
		ret[i] = envSample{
			InfluxMeasurement: "test",
			Time:              time.Now(),
			Location:          "Rm 243",
			Temperature:       70 + float64(i),
			Humidity:          60 - float64(i),
			ID:                "12432as32",
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
	c = c.UseDB(db)
	for _, p := range samples {
		err := c.WritePoint(p)
		if err != nil {
			log.Fatal("Error writing point: ", err)
		}
	}

	// query data from db
	samplesRead := []envSampleRead{}

	q := `SELECT * FROM test ORDER BY time DESC LIMIT 10`
	err = c.UseDB(db).DecodeQuery(q, &samplesRead)
	if err != nil {
		log.Fatal("Query error: ", err)
	}
	log.Printf("Samples read: %+v\n", samplesRead)
}
```

## Details

There are several advantages decoding and encoding data directly from Go
Structs:

1. The database bschema is documented by the Go type definition. This helps ensure
   data is written consistently to the database. When all your data is clearly
   defined in Go structs, it is much more obvious how to organize it, what goes
   in what measurement, when to create a new measurement, etc. When writing
   straight tags/values, it is much easier to create a disorganized mess.
1. All the code for decoding and encoding the various data types supported
   by InfluxDb are handled in one place, rather than repeating this logic over
   and over for every Query.
1. Likewise, code for handling arrays (translating Go array to InfluxDb fields
   like temp0, temp1, temp2, ...) can be in one place.
1. Reading and Writing data is much simpler and requires way less code.

Using Go reflection to automate data decode may be slightly slower
than custom decode logic for every query, but it seems the time to decode the
data will relatively fast compared to the time to run a InfluxDb query, so
may be negligable (this is an assumption at this point and has not been
proven).

The decode_test.go file contains a number of tests that illustrate the
conversion from influx JSON to Go struct values.

## Acknowledgments

The [mapstructure](https://github.com/mitchellh/mapstructure)
library provided a very useful reference for learning how to
use the Go reflect functionality.

## Status

Todo:

* [x] handle larger query datasets (multiple series, etc)
* [x] add write capability (directly write Go structs into influxdb)
* [x] add godoc documentation
* [x] get working with influxdb 1.7 client
* [ ] see if still applicable for influxdb 2.x
* [ ] decode/encode val0, val1, val2 fields in influx to Go array
* [ ] use Go struct field tags to help build SELECT statement
* [ ] optimize query for performace (pre-allocate slices, etc)
* [ ] come up with a better name (indecode, etc)
* [ ] finish error checking

Review/Pull requests welcome!

## License

MIT
