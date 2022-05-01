package middleware

import (
	"net/http"
	"saturday/src/util"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) == 0 {
		return
	}
	for _, err := range c.Errors {
		_, ok := util.IsServiceError(err)
		if !ok {
			util.Logger.Error(err)
		}
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": "Internal Server Error",
	})
}
