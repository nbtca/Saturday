package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBody(c *gin.Context, target interface{}) error {
	err := c.ShouldBindJSON(target)
	if err != nil {
		return err
	}
	ve := ValidationHandler(target)
	if ve == nil {
		return nil
	}
	serviceError := MakeServiceError(http.StatusUnprocessableEntity).SetMessage("Validation Failed")
	for _, fe := range ve {
		serviceError.AddDetailError(c.Request.URL.Path, fe.Field(), GetErrorMessage(fe))
	}
	return serviceError
}
