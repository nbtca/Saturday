package middleware

import (
	"log"
	"net/http"

	"github.com/nbtca/saturday/util"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) == 0 {
		return
	}
	for _, err := range c.Errors {
		log.Println(err)
		_, ok := util.IsServiceError(err)
		if !ok {
			util.LogErrorWithStackTrace("Internal server error occurred", err.Err)
		} else {
			util.Logger.Debugf("Service error occurred: %s", err.Error())
		}
	}
	c.JSON(util.MakeServiceError(http.StatusInternalServerError).
		SetMessage("Internal Server Error").
		Build())
}
