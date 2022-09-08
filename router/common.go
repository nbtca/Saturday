package router

import (
	"net/http"
	"saturday/model/dto"
	"saturday/util"

	"github.com/gin-gonic/gin"
)

type CommonRouter struct{}

func (CommonRouter) Upload(c *gin.Context) {
	maxBytes := 1024 * 1024 * 10 // 10MB

	file, err := c.FormFile("file")
	if err != nil {
		util.CheckError(c, util.MakeServiceError(http.StatusUnprocessableEntity).SetMessage(err.Error()))
		return
	}

	if file.Size > int64(maxBytes) {
		util.CheckError(c, util.MakeServiceError(http.StatusUnprocessableEntity).SetMessage("file size too large (max 10MB)"))
		return
	}

	reader, err := file.Open()
	name := file.Filename
	if err != nil {
		util.CheckError(c, util.MakeServiceError(http.StatusUnprocessableEntity).SetMessage(err.Error()))
		return
	}
	url, err := util.Upload(name, reader)
	if err != nil {
		util.CheckError(c, util.MakeServiceError(http.StatusUnprocessableEntity).SetMessage(err.Error()))
		return
	}
	res := dto.FileUploadResponse{
		Url: url,
	}
	c.JSON(200, res)
}

var CommonRouterApp = CommonRouter{}
