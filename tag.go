package influxdbhelper

import "strings"

func isInfluxTag(structTag string) bool {
	parts := strings.Split(structTag, ",")

	for _, part := range parts {
		if part == "tag" {
			return true
		}
	}

	return false
}

func getInfluxFieldTagName(structTag string) string {
	parts := strings.Split(structTag, ",")

	for _, part := range parts {
		if part != "tag" {
			return part
		}
	}

	return ""
}
