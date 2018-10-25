package influxdbhelper

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

func encode(d interface{}, timeField *usingValue) (t time.Time, tags map[string]string, fields map[string]interface{}, measurement string, err error) {
	tags = make(map[string]string)
	fields = make(map[string]interface{})
	dValue := reflect.ValueOf(d)

	if dValue.Kind() == reflect.Ptr {
		dValue = reflect.Indirect(dValue)
	}

	if dValue.Kind() != reflect.Struct {
		err = errors.New("data must be a struct")
		return
	}

	if timeField == nil {
		timeField = &usingValue{"time", false}
	}

	for i := 0; i < dValue.NumField(); i++ {
		f := dValue.Field(i)
		structFieldName := dValue.Type().Field(i).Name
		if structFieldName == "InfluxMeasurement" {
			measurement = f.String()
			continue
		}
		fieldTag := dValue.Type().Field(i).Tag.Get("influx")
		fieldData := getInfluxFieldTagData(structFieldName, fieldTag)

		if fieldData.fieldName == "-" {
			continue
		}

		if fieldData.fieldName == timeField.value {
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

	if measurement == "" {
		measurement = dValue.Type().Name()
	}

	return
}
