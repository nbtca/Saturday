package util

import (
	"reflect"
	"time"
)

func FieldsConstructor(q interface{}) []string {
	t := reflect.TypeOf(q)
	i := 0
	var res []string
	shouldAppendField := func(field reflect.StructField) (string, bool) {
		dbTag := t.Field(i).Tag.Get("json")
		// visibleTag := t.Field(i).Tag.Get("visible")
		if dbTag == "" {
			return "", false
		}
		return dbTag, true
	}
	if reflect.ValueOf(q).Kind() == reflect.Struct {
		for ; i < t.NumField()-1; i++ {
			if dbTag, should := shouldAppendField(t.Field(i)); should {
				res = append(res, dbTag)
			}
		}
	}
	if dbTag, should := shouldAppendField(t.Field(i)); should {
		res = append(res, dbTag)
	}
	return res
}

func GetDate() string {
	return time.Now().Format("2006-01-02 15:04:11")
}
