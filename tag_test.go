package influxdbhelper

import "testing"

func TestTag(t *testing.T) {
	data := []struct {
		fieldTag        string
		structFieldName string
		fieldName       string
		isTag           bool
		isField         bool
	}{
		{"", "Test", "Test", false, true},
		{"", "Test", "Test", false, true},
		{",tag", "Test", "Test", true, false},
		{",field,tag", "Test", "Test", true, true},
		{",tag,field", "Test", "Test", true, true},
		{",field", "Test", "Test", false, true},
		{"test", "Test", "test", false, true},
		{"test,tag", "Test", "test", true, false},
		{"test,field,tag", "Test", "test", true, true},
		{"test,tag,field", "Test", "test", true, true},
		{"test,field", "Test", "test", false, true},
	}

	for _, testData := range data {
		fieldData := getInfluxFieldTagData(testData.structFieldName, testData.fieldTag)
		if fieldData.fieldName != testData.fieldName {
			t.Errorf("%v != %v", fieldData.fieldName, testData.fieldName)
		}
		if fieldData.isField != testData.isField {
			t.Errorf("%v != %v", fieldData.isField, testData.isField)
		}
		if fieldData.isTag != testData.isTag {
			t.Errorf("%v != %v", fieldData.isTag, testData.isTag)
		}
	}
}
