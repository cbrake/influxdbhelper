package influxdbhelper

import (
	"reflect"
	"testing"

	"github.com/y0ssar1an/q"
)

func TestEncodeDataNotStruct(t *testing.T) {
	_, _, err := Encode([]int{1, 2, 3})
	if err == nil {
		t.Error("Expected error")
	}
}

func TestEncode(t *testing.T) {
	type MyType struct {
		TagValue     string  `influx:"tagValue,tag"`
		IntValue     int     `influx:"intValue"`
		FloatValue   float64 `influx:"floatValue"`
		BoolValue    bool    `influx:"boolValue"`
		StringValue  string  `influx:"stringValue"`
		IgnoredValue string  `influx:"-"`
	}

	d := MyType{
		"tag-value",
		10,
		10.5,
		true,
		"string",
		"ignored",
	}

	tagsExp := map[string]string{
		"tagValue": "tag-value",
	}

	fieldsExp := map[string]interface{}{
		"intValue":    d.IntValue,
		"floatValue":  d.FloatValue,
		"boolValue":   d.BoolValue,
		"stringValue": d.StringValue,
	}

	tags, fields, err := Encode(d)

	if err != nil {
		t.Error("Error encoding: ", err)
	}

	if !reflect.DeepEqual(tags, tagsExp) {
		t.Error("tags not encoded correctly")
	}

	if !reflect.DeepEqual(fields, fieldsExp) {
		q.Q(fields)
		q.Q(fieldsExp)
		t.Error("fields not encoded correctly")
	}
}
