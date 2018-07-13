package influxdbhelper

import (
	"regexp"
	"strings"

	client "github.com/influxdata/influxdb/client/v2"
)

var reRemoveExtraSpace = regexp.MustCompile(`\s\s+`)

// CleanQuery can be used to strip a query string of
// newline characters. Typically only used for debugging.
func CleanQuery(query string) string {
	ret := strings.Replace(query, "\n", "", -1)
	ret = reRemoveExtraSpace.ReplaceAllString(ret, " ")
	return ret
}

// A Client represents an influxdbhelper client connection to
// an InfluxDb server.
type Client struct {
	url       string
	client    client.Client
	precision string
}

// NewClient returns a new influxdbhelper client given a url, user,
// password, and precision strings.
//
// url is typically something like: http://localhost:8086
//
// precision can be ‘h’, ‘m’, ‘s’, ‘ms’, ‘u’, or ‘ns’ and is
// used during write operations.
func NewClient(url, user, passwd, precision string) (*Client, error) {
	ret := Client{
		url:       url,
		precision: precision,
	}

	client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     url,
		Username: user,
		Password: passwd,
	})

	ret.client = client

	return &ret, err
}

// InfluxClient returns the influxdb/client/v2 client if low level
// queries or writes need to be executed.
func (c Client) InfluxClient() client.Client {
	return c.client
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
func (c Client) Query(db, cmd string, result interface{}) (err error) {
	query := client.Query{
		Command:   cmd,
		Database:  db,
		Chunked:   false,
		ChunkSize: 100,
	}

	var response *client.Response
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

	err = decode(series.Columns, series.Values, result)

	return
}

// WritePoint is used to write arbitrary data into InfluxDb.
//
// data must be a struct with struct field tags that defines the names used
// in InfluxDb for each field. A "tag" tag can be added to indicate the
// struct field should be an InfluxDb tag (vs field). A tag of '-' indicates
// the struct field should be ignored. A struct field of Time is required and
// is used for the time of the sample.
func (c Client) WritePoint(db, measurement string, data interface{}) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: c.precision,
	})

	if err != nil {
		return err
	}

	t, tags, fields, err := encode(data)

	if err != nil {
		return err
	}

	pt, err := client.NewPoint(measurement, tags, fields, t)

	if err != nil {
		return err
	}

	bp.AddPoint(pt)

	return c.client.Write(bp)
}
