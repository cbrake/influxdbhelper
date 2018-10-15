package influxdbhelper

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	influxClient "github.com/influxdata/influxdb/client/v2"
)

var reRemoveExtraSpace = regexp.MustCompile(`\s\s+`)

// CleanQuery can be used to strip a query string of
// newline characters. Typically only used for debugging.
func CleanQuery(query string) string {
	ret := strings.Replace(query, "\n", "", -1)
	ret = reRemoveExtraSpace.ReplaceAllString(ret, " ")
	return ret
}

type Client interface {
	influxClient.Client

	// UseDB sets the DB to use for Query, WritePoint, and WritePointTagsFields
	UseDB(db string) Client

	// UseMeasurement sets the measurment to use for WritePoint, and WritePointTagsFields
	UseMeasurement(measurement string) Client

	// Query executes an InfluxDb query, and unpacks the result into the
	// result data structure.
	DecodeQuery(query string, result interface{}) error

	// WritePoint is used to write arbitrary data into InfluxDb.
	WritePoint(data interface{}) error

	// WritePointTagsFields is used to write a point specifying tags and fields.
	WritePointTagsFields(tags map[string]string, fields map[string]interface{}, t time.Time) error
}

// A Client represents an influxdbhelper influxClient connection to
// an InfluxDb server.
type helperClient struct {
	url       string
	client    influxClient.Client
	precision string
	using     *helperUsing
}

type helperUsing struct {
	db string
	measurement string
}

// NewClient returns a new influxdbhelper influxClient given a url, user,
// password, and precision strings.
//
// url is typically something like: http://localhost:8086
//
// precision can be ‘h’, ‘m’, ‘s’, ‘ms’, ‘u’, or ‘ns’ and is
// used during write operations.
func NewClient(url, user, passwd, precision string) (Client, error) {
	ret := &helperClient{
		url:       url,
		precision: precision,
	}

	client, err := influxClient.NewHTTPClient(influxClient.HTTPConfig{
		Addr:     url,
		Username: user,
		Password: passwd,
	})

	ret.client = client

	return ret, err
}

// Ping checks that status of cluster, and will always return 0 time and no
// error for UDP clients.
func (c *helperClient) Ping(timeout time.Duration) (time.Duration, string, error) {
	return c.client.Ping(timeout)
}

// Write takes a BatchPoints object and writes all Points to InfluxDB.
func (c *helperClient) Write(bp influxClient.BatchPoints) error {
	return c.client.Write(bp)
}

// Query makes an InfluxDB Query on the database. This will fail if using
// the UDP client.
func (c *helperClient) Query(q influxClient.Query) (*influxClient.Response, error) {
	return c.client.Query(q)
}

// Close releases any resources a Client may be using.
func (c *helperClient) Close() error {
	return c.client.Close()
}

// UseDB sets the DB to use for Query, WritePoint, and WritePointTagsFields
func (c *helperClient) UseDB(db string) Client {
	if c.using == nil {
		c.using = &helperUsing{}
	}

	c.using.db = db
	return c
}

// UseMeasurement sets the DB to use for Query, WritePoint, and WritePointTagsFields
func (c *helperClient) UseMeasurement(measurement string) Client {
	if c.using == nil {
		c.using = &helperUsing{}
	}

	c.using.measurement = measurement
	return c
}


// Query executes an InfluxDb query, and unpacks the result into the
// result data structure.
//
// result must be an array of structs that contains the fields returned
// by the query. The struct type must always contain a Time field. The
// struct type must also include influx field tags which map the struct
// field name to the InfluxDb field/tag names. This tag is currently
// required as typically Go struct field names start with a capital letter,
// and InfluxDb field/tag names typically start with a lower case letter.
// The struct field tag can be set to '-' which indicates this field
// should be ignored.
func (c *helperClient) DecodeQuery(q string, result interface{}) (err error) {
	if c.using == nil {
		return fmt.Errorf("no db set for query")
	}

	query := influxClient.Query{
		Command:   q,
		Database:  c.using.db,
		Chunked:   false,
		ChunkSize: 100,
	}

	var response *influxClient.Response
	response, err = c.client.Query(query)

	if response.Error() != nil {
		return response.Error()
	}

	if err != nil {
		return
	}

	results := response.Results
	if len(results) < 1 || len(results[0].Series) < 1 {
		return
	}

	series := results[0].Series[0]
	err = decode(series, result)

	return
}

// WritePoint is used to write arbitrary data into InfluxDb.
//
// data must be a struct with struct field tags that defines the names used
// in InfluxDb for each field. A "tag" tag can be added to indicate the
// struct field should be an InfluxDb tag (vs field). A tag of '-' indicates
// the struct field should be ignored. A struct field of Time is required and
// is used for the time of the sample.
func (c *helperClient) WritePoint(data interface{}) error {
	if c.using == nil {
		return fmt.Errorf("no db set for query")
	}

	t, tags, fields, measurement, err := encode(data)

	if c.using.measurement == "" {
		c.using.measurement = measurement
	}

	if err != nil {
		return err
	}

	return c.WritePointTagsFields( tags, fields, t)
}


// WritePointTagsFields is used to write a point specifying tags and fields.
func (c *helperClient) WritePointTagsFields(tags map[string]string, fields map[string]interface{}, t time.Time) (err error) {
	if c.using == nil {
		return fmt.Errorf("no db set for query")
	}

	bp, err := influxClient.NewBatchPoints(influxClient.BatchPointsConfig{
		Database:  c.using.db,
		Precision: c.precision,
	})

	if err != nil {
		return err
	}

	pt, err := influxClient.NewPoint(c.using.measurement, tags, fields, t)

	if err != nil {
		return err
	}

	bp.AddPoint(pt)

	return c.client.Write(bp)
}
