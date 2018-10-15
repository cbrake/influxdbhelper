package influxdbhelper

import (
	"reflect"
	"testing"
	"time"
)

func TestEncodeDataNotStruct(t *testing.T) {
	_, _, _, _, err := encode([]int{1, 2, 3}, nil)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestEncodeSetsMesurment(t *testing.T) {
	type MyType struct {
		Val string `influx:"val"`
	}

	d := &MyType{"test-data"}
	_, _, _, measurement, err := encode(d, nil)

	if err != nil {
		t.Error("Error encoding: ", err)
	}

	if measurement != "MyType" {
		t.Errorf("%v != %v", measurement, "MyType")
	}
}

func TestEncodeUsesTimeField(t *testing.T) {
	type MyType struct {
		MyTimeField             time.Time `influx:"my_time_field"`
		Val string `influx:"val"`
	}

	td, _ := time.Parse(time.RFC822, "27 Oct 78 15:04 PST")

	d := &MyType{td,"test-data"}
	tv, _, _, _, err := encode(d, &usingValue{"my_time_field", false})

	if tv != td {
		t.Error("Did not properly use the time field specified")
	}

	if err != nil {
		t.Error("Error encoding: ", err)
	}
}

func TestEncode(t *testing.T) {
	type MyType struct {
		InfluxMeasurement      Measurement
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
		"test",
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

	tm, tags, fields, measurement, err := encode(d, nil)

	if err != nil {
		t.Error("Error encoding: ", err)
	}

	if measurement != d.InfluxMeasurement {
		t.Errorf("%v != %v", measurement, d.InfluxMeasurement)
	}

	if _, ok := fields["InfluxMeasurement"]; ok {
		t.Errorf("Found InfluxMeasurement in the fields!")
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
