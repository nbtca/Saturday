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
	Body       T
	Total      *int64  `header:"X-Total-Count" doc:"Total number of items"`
	Offset     *uint64 `header:"X-Offset" doc:"Offset for pagination"`
	Limit      *uint64 `header:"X-Limit" doc:"Limit for pagination"`
	Page       *int64  `header:"X-Page" doc:"Current page number (1-indexed)"`
	TotalPages *int64  `header:"X-Total-Pages" doc:"Total number of pages"`
}

func MakeCommonResponse[T any](body T) *CommonResponse[T] {
	return &CommonResponse[T]{Body: body}
}

func MakePaginatedResponse[T any](body T, total int64, offset, limit uint64) *CommonResponse[T] {
	page := int64(1)
	if limit > 0 {
		page = int64(offset/limit) + 1
	}
	totalPages := int64(0)
	if limit > 0 {
		totalPages = (total + int64(limit) - 1) / int64(limit)
	}
	return &CommonResponse[T]{
		Body:       body,
		Total:      &total,
		Offset:     &offset,
		Limit:      &limit,
		Page:       &page,
		TotalPages: &totalPages,
	}
}
