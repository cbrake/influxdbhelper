package influxdbhelper

import (
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"testing"
	"time"

	influxModels "github.com/influxdata/influxdb1-client/models"
)

func TestDecode(t *testing.T) {
	data := influxModels.Row{
		Name: "bla",
		Columns: []string{
			"intValue",
			"floatValue",
			"boolValue",
			"stringValue",
		},
		Values: make([][]interface{}, 0),
		Tags:   map[string]string{"tagValue": "tag-value"},
	}

	type DecodeType struct {
		TagValue     string  `influx:"tagValue,tag"`
		IntValue     int     `influx:"intValue"`
		FloatValue   float64 `influx:"floatValue"`
		BoolValue    bool    `influx:"boolValue"`
		StringValue  string  `influx:"stringValue"`
		IgnoredValue string  `influx:"-"`
	}

	expected := []DecodeType{}

	for i := 0; i < 10; i++ {
		v := DecodeType{
			"tag-value",
			i,
			float64(i),
			math.Mod(float64(i), 2) == 0,
			strconv.Itoa(i),
			"",
		}

		vI := []interface{}{
			v.IntValue,
			v.FloatValue,
			v.BoolValue,
			v.StringValue,
		}

		expected = append(expected, v)
		data.Values = append(data.Values, vI)

	}

	decoded := []DecodeType{}

	err := decode([]influxModels.Row{data}, &decoded)
	if err != nil {
		t.Error("Error decoding: ", err)
	}

	if !reflect.DeepEqual(expected, decoded) {
		t.Error("decoded value is not right")
	}
}

func TestDecodeMissingColumn(t *testing.T) {
	data := influxModels.Row{
		Name: "bla",
		Columns: []string{
			"val1",
		},
		Values: make([][]interface{}, 0),
		Tags:   map[string]string{},
	}

	type DecodeType struct {
		Val1 int `influx:"val1"`
		Val2 int `influx:"val2"`
	}

	expected := []DecodeType{{1, 0}}
	data.Values = append(data.Values, []interface{}{1})
	decoded := []DecodeType{}
	err := decode([]influxModels.Row{data}, &decoded)

	if err != nil {
		t.Error("UnExpected error decoding: ", data, &decoded)
	}

	if !reflect.DeepEqual(expected, decoded) {
		t.Error("decoded value is not right")
	}
}

func TestDecodeWrongType(t *testing.T) {
	data := influxModels.Row{
		Name: "bla",
		Columns: []string{
			"val1",
			"val2",
		},
		Values: make([][]interface{}, 0),
		Tags:   map[string]string{},
	}

	type DecodeType struct {
		Val1 int     `influx:"val1"`
		Val2 float64 `influx:"val2"`
	}

	expected := []DecodeType{{1, 2.0}}
	data.Values = append(data.Values, []interface{}{1.0, 2})
	decoded := []DecodeType{}
	err := decode([]influxModels.Row{data}, &decoded)
	if err != nil {
		t.Error("Unexpected error decoding: ", err, data, decoded)
	}

	if !reflect.DeepEqual(expected, decoded) {
		t.Error("decoded value is not right", expected, decoded)
	}
}

func TestDecodeTime(t *testing.T) {
	data := influxModels.Row{
		Name: "bla",
		Columns: []string{
			"time",
			"value",
		},
		Values: make([][]interface{}, 0),
		Tags:   map[string]string{},
	}

	type DecodeType struct {
		Time  time.Time `influx:"time"`
		Value float64   `influx:"value"`
	}

	timeS := "2018-06-14T21:47:11Z"
	ti, err := time.Parse(time.RFC3339, timeS)
	if err != nil {
		t.Error("error parsing expected time: ", err)
	}

	expected := []DecodeType{{ti, 2.0}}
	data.Values = append(data.Values, []interface{}{timeS, 2.0})
	decoded := []DecodeType{}
	err = decode([]influxModels.Row{data}, &decoded)

	if err != nil {
		t.Error("Error decoding: ", err)
	}

	if !reflect.DeepEqual(expected, decoded) {
		t.Error("decoded value is not right")
	}
}

func TestDecodeJsonNumber(t *testing.T) {
	data := influxModels.Row{
		Name: "bla",
		Columns: []string{
			"val1",
			"val2",
		},
		Values: make([][]interface{}, 0),
		Tags:   map[string]string{},
	}

	type DecodeType struct {
		Val1 int     `influx:"val1"`
		Val2 float64 `influx:"val2"`
	}

	expected := []DecodeType{{1, 2.0}}
	data.Values = append(data.Values, []interface{}{json.Number("1"), json.Number("2.0")})
	decoded := []DecodeType{}
	err := decode([]influxModels.Row{data}, &decoded)

	if err != nil {
		t.Error("Error decoding: ", err)
	}

	if !reflect.DeepEqual(expected, decoded) {
		t.Error("decoded value is not right")
	}
}

func TestDecodeUnsedStructValue(t *testing.T) {
	data := influxModels.Row{
		Name: "bla",
		Columns: []string{
			"val1",
			"val2",
		},
		Values: make([][]interface{}, 0),
		Tags:   map[string]string{},
	}

	type DecodeType struct {
		Val1 int     `influx:"val1"`
		Val2 float64 `influx:"-"`
	}

	expected := []DecodeType{{1, 0}}
	data.Values = append(data.Values, []interface{}{1, 1.1})
	decoded := []DecodeType{}
	err := decode([]influxModels.Row{data}, &decoded)

	if err != nil {
		t.Error("Error decoding: ", err)
	}

	if !reflect.DeepEqual(expected, decoded) {
		t.Error("decoded value is not right")
	}
}

func TestDecodeMeasure(t *testing.T) {
	data := influxModels.Row{
		Name: "bla",
		Columns: []string{
			"val1",
			"val2",
		},
		Values: make([][]interface{}, 0),
		Tags:   map[string]string{},
	}

	type DecodeType struct {
		InfluxMeasurement Measurement
		Val1              int     `influx:"val1"`
		Val2              float64 `influx:"-"`
	}

	expected := []DecodeType{{"bla", 1, 0}}
	data.Values = append(data.Values, []interface{}{1, 1.1})
	decoded := []DecodeType{}
	err := decode([]influxModels.Row{data}, &decoded)

	if decoded[0].InfluxMeasurement != expected[0].InfluxMeasurement {
		t.Error("Decoded Wrong measure")
	}

	if err != nil {
		t.Error("Error decoding: ", err)
	}

	if !reflect.DeepEqual(expected, decoded) {
		t.Error("decoded value is not right")
	}
}
