package influxdbhelper

import "strings"

type influxFieldTagData struct {
	fieldName string
	isTag     bool
	isField   bool
}

func getInfluxFieldTagData(fieldName, structTag string) (fieldData *influxFieldTagData) {
	fieldData = &influxFieldTagData{fieldName: fieldName}
	parts := strings.Split(structTag, ",")
	fieldName, parts = parts[0], parts[1:]
	if fieldName != "" {
		fieldData.fieldName = fieldName
	}

	for _, part := range parts {
		if part == "tag" {
			fieldData.isTag = true
		}
		if part == "field" {
			fieldData.isField = true
		}
	}

	if !fieldData.isField && !fieldData.isTag {
		fieldData.isField = true
	}

	return
}
