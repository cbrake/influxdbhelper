package influxdbhelper

import (
	"reflect"
	"time"

	influxModels "github.com/influxdata/influxdb/models"
	"github.com/mitchellh/mapstructure"
)



// Decode is used to process data returned by an InfluxDb query and uses reflection
// to transform it into an array of structs of type result.
//
// This function is used internally by the Query function.
func decode(influxResult influxModels.Row, result interface{}) error {
	influxData := make([]map[string]interface{}, 0)

	for _, v := range influxResult.Values {
		r := make(map[string]interface{})
		for i, c := range influxResult.Columns {
			if len(v) >= i+1 {
				r[c] = v[i]
			}
		}
		for tag, val := range influxResult.Tags {
			r[tag] = val
		}
		r["InfluxMeasurement"] = influxResult.Name

		influxData = append(influxData, r)
	}

	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result: result,
		TagName: "influx",
		WeaklyTypedInput: false,
		ZeroFields: false,
		DecodeHook: func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			if t == reflect.TypeOf(time.Time{}) && f == reflect.TypeOf("") {
				return time.Parse(time.RFC3339, data.(string))
			}

			return data, nil
		},
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(influxData)
}
