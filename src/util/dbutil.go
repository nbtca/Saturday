package util

import (
	"reflect"
)

func FieldsConstructor(q interface{}) string {
	t := reflect.TypeOf(q)
	i := 0
	res := ""
	if reflect.ValueOf(q).Kind() == reflect.Struct {
		for ; i < t.NumField()-1; i++ {
			dbTag := t.Field(i).Tag.Get("json")
			if dbTag != "" {
				res += dbTag + ", "
			}
		}
	}
	dbTag := t.Field(i).Tag.Get("json")
	if dbTag != "" {
		res += dbTag
	}
	return res
}
