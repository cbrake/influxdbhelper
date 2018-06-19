package influxdbhelper

import (
	"errors"

	client "github.com/influxdata/influxdb/client/v2"
)

type Client struct {
	url    string
	client client.Client
}

func NewClient(url, user, passwd string) (*Client, error) {
	ret := Client{
		url: url,
	}

	client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     url,
		Username: user,
		Password: passwd,
	})

	ret.client = client

	return &ret, err
}

func (h Client) Query(db, cmd string, result interface{}) (err error) {
	query := client.Query{
		Command:   cmd,
		Database:  db,
		Chunked:   true,
		ChunkSize: 100,
	}

	var response *client.Response
	response, err = h.client.Query(query)

	if response.Error() != nil {
		return response.Error()
	}

	if err != nil {
		return
	}

	results := response.Results
	if len(results) < 1 || len(results[0].Series) < 1 {
		err = errors.New("No data returned")
		return
	}

	series := results[0].Series[0]

	err = Decode(series.Columns, series.Values, result)

	return
}
