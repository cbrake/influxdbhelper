package influxdbhelper

import (
	"reflect"
	"testing"
	"time"
)

func TestEncodeDataNotStruct(t *testing.T) {
	_, _, _, err := encode([]int{1, 2, 3})
	if err == nil {
		t.Error("Expected error")
	}
}

func TestEncode(t *testing.T) {
	type MyType struct {
		Time         time.Time `influx:"time"`
		TagValue     string    `influx:"tagValue,tag"`
		IntValue     int       `influx:"intValue"`
		FloatValue   float64   `influx:"floatValue"`
		BoolValue    bool      `influx:"boolValue"`
		StringValue  string    `influx:"stringValue"`
		IgnoredValue string    `influx:"-"`
	}

	d := MyType{
		time.Now(),
		"tag-value",
		10,
		10.5,
		true,
		"string",
		"ignored",
	}

	timeExp := d.Time

	tagsExp := map[string]string{
		"tagValue": "tag-value",
	}

	fieldsExp := map[string]interface{}{
		"intValue":    d.IntValue,
		"floatValue":  d.FloatValue,
		"boolValue":   d.BoolValue,
		"stringValue": d.StringValue,
	}

	tm, tags, fields, err := encode(d)

	if err != nil {
		t.Error("Error encoding: ", err)
	}

	if !tm.Equal(timeExp) {
		t.Error("Time does not match")
	}

	if !reflect.DeepEqual(tags, tagsExp) {
		t.Error("tags not encoded correctly")
	}

	if !reflect.DeepEqual(fields, fieldsExp) {
		t.Error("fields not encoded correctly")
	}
}
