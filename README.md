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

type PumpEvent struct {
	Time      time.Time     `influx:"time"`
	Duration  time.Duration `influx:"-"`
	DurationS float64       `influx:"durationS"`
        PumpIndex string        `influx:"pumpIndex"`
        Value     float64       `influx:"value"`
}

query := `SELECT "durationS","pumpIndex","value"
	from myMeasurement
	order by time desc
	limit 50`

var events []PumpEvent

err = client.Query("mydb", query, &events)
```

## Details

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

* [ ] handle larger query datasets (multiple series, etc)
* [ ] add write capability (directly write Go structs into influxdb)
* [ ] use Go struct field tags to help build SELECT statement
* [ ] optimize query for performace (pre-allocate slices, etc)
* [ ] decode/encode val0, val1, val2 fields in influx to Go array

Pull requests welcome!

## License

MIT
