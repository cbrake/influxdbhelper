package influxdbhelper

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

func encode(d interface{}) (t time.Time, tags map[string]string, fields map[string]interface{}, err error) {
	tags = make(map[string]string)
	fields = make(map[string]interface{})
	dValue := reflect.ValueOf(d)
	if dValue.Kind() != reflect.Struct {
		err = errors.New("data must be a struct")
		return
	}

	for i := 0; i < dValue.NumField(); i++ {
		f := dValue.Field(i)
		fieldName := dValue.Type().Field(i).Name
		fieldTag := dValue.Type().Field(i).Tag.Get("influx")
		fieldData := getInfluxFieldTagData(fieldName, fieldTag)

		if fieldData.fieldName == "-" {
			continue
		}

		if fieldData.fieldName == "time" {
			// TODO error checking
			t = f.Interface().(time.Time)
			continue
		}

		if fieldData.isTag {
			tags[fieldData.fieldName] = fmt.Sprintf("%v", f)
		}

		if fieldData.isField {
			fields[fieldData.fieldName] = f.Interface()
		}
	}

	return
}
