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
		Time             time.Time `influx:"time"`
		TagValue         string    `influx:"tagValue,tag"`
		TagAndFieldValue string    `influx:"tagAndFieldValue,tag,field"`
		IntValue         int       `influx:"intValue"`
		FloatValue       float64   `influx:"floatValue"`
		BoolValue        bool      `influx:"boolValue"`
		StringValue      string    `influx:"stringValue"`
		StructFieldName  string    `influx:""`
		IgnoredValue     string    `influx:"-"`
	}

	d := MyType{
		time.Now(),
		"tag-value",
		"tag-and-field-value",
		10,
		10.5,
		true,
		"string",
		"struct-field",
		"ignored",
	}

	timeExp := d.Time

	tagsExp := map[string]string{
		"tagValue":         "tag-value",
		"tagAndFieldValue": "tag-and-field-value",
	}

	fieldsExp := map[string]interface{}{
		"tagAndFieldValue": d.TagAndFieldValue,
		"intValue":         d.IntValue,
		"floatValue":       d.FloatValue,
		"boolValue":        d.BoolValue,
		"stringValue":      d.StringValue,
		"StructFieldName":  d.StructFieldName,
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
