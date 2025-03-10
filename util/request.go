package util

import (
	"net/http"
	"strconv"

	"github.com/nbtca/saturday/model"

	"github.com/gin-gonic/gin"
)

func BindAll(c *gin.Context, target interface{}) error {
	c.ShouldBindUri(target)
	jsonBindingErr := c.ShouldBindJSON(target)
	if jsonBindingErr != nil {
		if jsonBindingErr.Error() == "unexpected EOF" {
			return MakeServiceError(http.StatusBadRequest).
				SetMessage("Problems parsing JSON")
		}
	}
	err := c.ShouldBindQuery(target)
	if err != nil {
		Logger.Print(err)
		return MakeValidationError(c.Request.URL.Path, err)
	}
	return nil
}

func GetPaginationQuery(c *gin.Context) (offset uint64, limit uint64, err error) {
	offset, err = strconv.ParseUint(c.DefaultQuery("offset", "0"), 10, 64)
	if err != nil {
		return
	}
	limit, err = strconv.ParseUint(c.DefaultQuery("limit", "50"), 10, 64)
	if err != nil {
		return
	}
	return
}

func GetIdentity(c *gin.Context) model.Identity {
	id := c.GetString("id")
	rawMember, _ := c.Get("member")
	member := rawMember.(model.Member)
	role := c.GetString("role")
	return model.Identity{
		Id:     id,
		Member: member,
		Role:   role,
	}
}

type CommonResponse[T any] struct {
	Body T
}

func MakeCommonResponse[T any](body T) *CommonResponse[T] {
	return &CommonResponse[T]{Body: body}
}
