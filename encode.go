package influxdbhelper

import (
	"errors"
	"reflect"
)

func Encode(d interface{}) (tags map[string]string,
	fields map[string]interface{}, err error) {
	tags = make(map[string]string)
	fields = make(map[string]interface{})
	dValue := reflect.ValueOf(d)
	if dValue.Kind() != reflect.Struct {
		err = errors.New("data must be a struct")
		return
	}

	for i := 0; i < dValue.NumField(); i++ {
		f := dValue.Field(i)
		fieldTag := dValue.Type().Field(i).Tag.Get("influx")

		isTag := isInfluxTag(fieldTag)
		name := getInfluxFieldTagName(fieldTag)

		if name == "-" {
			continue
		}

		if isTag {
			tags[name] = f.String()
		} else {
			fields[name] = f.Interface()
		}
	}

	return
}
