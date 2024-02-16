package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nbtca/saturday/util"
	"github.com/sirupsen/logrus"
)

func Logger(c *gin.Context) {
	startTime := time.Now()

	c.Next()

	endTime := time.Now()

	latencyTime := endTime.Sub(startTime)

	log.Println(c.GetTime("time_stamp").Format(time.RFC1123))
	util.Logger.WithFields(logrus.Fields{
		"status_code":   c.Writer.Status(),
		"latency":       latencyTime,
		"client_ip":     c.ClientIP(),
		"time_stamp":    c.GetTime("time_stamp").Format(time.RFC1123),
		"method":        c.Request.Method,
		"path":          c.Request.URL.Path,
		"request_proto": c.Request.Proto,
		"user_agent":    c.Request.UserAgent(),
		"error_message": c.Errors.String(),
		"userId":        c.GetString("id"),
	}).Info()

}
