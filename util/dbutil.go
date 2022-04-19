package util

import (
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
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

func GetPaginationQuery(c *gin.Context) (offset uint64, limit uint64, err error) {
	offset, err = strconv.ParseUint(c.DefaultQuery("offset", "0"), 10, 64)
	if err != nil {
		return
	}
	limit, err = strconv.ParseUint(c.DefaultQuery("offset", "30"), 10, 64)
	if err != nil {
		return
	}
	return
}
