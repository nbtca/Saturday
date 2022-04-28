package util

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BindAll(c *gin.Context, target interface{}) error {
	c.ShouldBindUri(target)
	c.ShouldBindJSON(target)
	err := c.ShouldBindQuery(target)
	if err != nil {
		return MakeValidationError(c.Request.URL.Path, err)
	}
	return nil
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

func SwapObject(from interface{}, to interface{}) error {
	dataByte, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(dataByte, to)
}
