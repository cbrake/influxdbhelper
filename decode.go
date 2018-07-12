package influxdbhelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func Decode(columns []string, values [][]interface{}, result interface{}) error {
	colIndex := map[string]int{}
	for i, col := range columns {
		colIndex[col] = i
	}

	resultV := reflect.ValueOf(result)
	if resultV.Kind() != reflect.Ptr {
		return errors.New("result must be ptr")
	}

	resultSlice := resultV.Elem()

	if !resultSlice.CanAddr() {
		return errors.New("result must be addressable (a pointer)")
	}

	if resultSlice.Kind() != reflect.Slice {
		return errors.New("result must be ptr to slice")
	}

	resultStruct := resultSlice.Type().Elem()
	if resultStruct.Kind() != reflect.Struct {
		return errors.New("result must be slice of structs")
	}

	numFields := resultStruct.NumField()
	resultStructFields := []reflect.StructField{}
	resultStructTags := []string{}

	for i := 0; i < numFields; i++ {
		f := resultStruct.Field(i)
		resultStructFields = append(resultStructFields,
			f)
		resultStructTags = append(resultStructTags,
			f.Tag.Get("influx"))
	}

	// not sure why we need to do this, but we need to Set resultSlice
	// at the end of this function for things to work
	resultSlice_ := resultSlice

	// Accumulate any errors
	errs := make([]string, 0)

	missingTagError := false
	typeError := false

	for _, vIn := range values {
		vOut := reflect.Indirect(reflect.New(resultStruct))
		valueCount := 0
		for i := 0; i < vOut.NumField(); i++ {
			f := vOut.Field(i)
			// FIXME, not sure how to get the tags
			// from vOut
			tag := resultStructTags[i]
			if tag == "-" {
				continue
			}

			tag = getInfluxFieldTagName(tag)
			i, ok := colIndex[tag]

			if !ok {
				if !missingTagError {
					err := errors.New(fmt.Sprintf("Missing tag: %v", tag))
					errs = AppendErrors(errs, err)
					missingTagError = true
				}
				continue
			}

			if vIn[i] == nil {
				continue
			}

			if f.Type() == reflect.TypeOf(time.Time{}) {
				timeS, ok := vIn[i].(string)
				if !ok {
					e := errors.New("Time input is not string")
					errs = AppendErrors(errs, e)
				} else {
					time, err := time.Parse(time.RFC3339, timeS)
					if err != nil {
						e := errors.New("Error parsing time")
						errs = AppendErrors(errs, e)
					} else {
						vIn[i] = time
					}
				}
			}

			if reflect.TypeOf(vIn[i]) == reflect.TypeOf(json.Number("1")) {
				if f.Type() == reflect.TypeOf(1.0) {
					vIn_, _ := vIn[i].(json.Number)
					vIn__, err := strconv.ParseFloat(string(vIn_), 64)
					if err != nil {
						es := "error converting json.Number"
						errs = AppendErrors(errs, errors.New(es))
					}
					vIn[i] = vIn__
				} else {
					vIn_, _ := vIn[i].(json.Number)
					vIn__, err := strconv.Atoi(string(vIn_))
					if err != nil {
						es := "error converting json.Number"
						errs = AppendErrors(errs, errors.New(es))
					}
					vIn[i] = vIn__
				}
			}

			if reflect.TypeOf(vIn[i]) != f.Type() {
				if !typeError {
					es := fmt.Sprintf("Type mismatch on decode of %v: %v != %v",
						vIn[i],
						reflect.TypeOf(vIn[i]).String(),
						f.Type().String())

					errs = AppendErrors(errs, errors.New(es))
					typeError = true
				}
				continue
			}
			f.Set(reflect.ValueOf(vIn[i]))
			valueCount = 1
		}

		if valueCount > 0 {
			resultSlice_ = reflect.Append(resultSlice_, vOut)
		}
	}

	resultSlice.Set(resultSlice_)

	if len(errs) > 0 {
		return &Error{errs}
	}

	return nil
}
