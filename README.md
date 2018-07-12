# InfluxDb Helper Library

> easily write and query InfluxDb from Go programs

This library allows you to encode/decode InfluxDb data to/from
Go structs -- similiar to JSON and MongoDb using Go struct field tags.

## Install

```
go get github.com/influxdata/influxdb/client/v2
go get github.com/cbrake/influxdbhelper
```

## Example

```go
client, err = influxdbhelper.NewClient("http://localhost:8086", "user", "passwd")

if err != nil {
	...
}

// structs used to read/write influxdb must contain a Time field,
// and all fields marked with the "tag" tag must be a string, as InfluxDb
// only allows strings to be used as tags.
type PumpEvent struct {
	Time      time.Time     `influx:"time"`
	Duration  time.Duration `influx:"-"`              // ignored field
	DurationS float64       `influx:"durationS"`
        PumpIndex string        `influx:"pumpIndex,tag"`  // written as tag (vs field)
        Value     float64       `influx:"value"`
}

p := PumpEvent{
	Time: time.Now(),
	Duration: time.Minute*2,
	DurationS: 60*2,
	PumpIndex: "1",
	Value: 350,
}

// WritePoint uses PumpIndex as an InfluxDb tag, and the rest of the struct fields as
// InfluxDb fields.
err = client.WritePoint(db, "myMeasurement", p)

query := `SELECT "durationS","pumpIndex","value"
	from myMeasurement
	order by time desc
	limit 50`

var events []PumpEvent

// Query populates the events slice with the data returned from the
// influxdb query.
err = client.Query("mydb", query, &events)
```

## Details

There are several advantages decoding and encoding data directly from Go
Structs:

1. The data schema is documented by the Go type definition. This helps ensure
   data is written consistently to the database. When all your data is clearly
   defined in Go structs, it is much more obvious how to organize it, and how to
   organize data in measurements. When writing straight tags/values, it is
   much easier to create a disorganized mess.
1. All the code for decoding and encoding the various data types supported
   by InfluxDb are handled in one place, rather than repeating this logic over
   and over for every Query.
1. Likewise, code for handling arrays can be in one place.
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

This library is currently in the proof of concept phase, and the code is not
optimized for performance, nor is it very clean at this point. If there are other
libraries that do similiar things, I would be very interested in learning about them.

Todo:

* [x] handle larger query datasets (multiple series, etc)
* [x] add write capability (directly write Go structs into influxdb)
* [ ] use Go struct field tags to help build SELECT statement
* [ ] optimize query for performace (pre-allocate slices, etc)
* [ ] decode/encode val0, val1, val2 fields in influx to Go array
* [ ] come up with a better name (indecode, ingodec, etc)
* [ ] add godoc documentation

Pull requests welcome!

## License

MIT
