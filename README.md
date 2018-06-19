# InfluxDb Helper Library

> easily write and query influxdb from Go programs

This library allows you to interface with InfluxDb databases with
Go structs -- similiar to how you are using to interfacing with
JSON data, MongoDb, and other databases. Struct field tags are
used to map InfluxDb tags and values to struct fields.

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

# Acknowledgments

The [mapstructure](https://github.com/mitchellh/mapstructure)
library provided a very useful reference for learning how to
use the Go reflect functionality.

# Status

This library is currently in the proof of concept phase, and the code is not
optimized for performance, nor is it very clean at this point.

Todo:

* [ ] handle larger query datasets (multiple series, etc)
* [ ] add write capability (directly write Go structs into influxdb)
* [ ] use Go struct field tags to help build SELECT statement
* [ ] optimize query for performace (pre-allocate slices, etc)

# License

MIT
