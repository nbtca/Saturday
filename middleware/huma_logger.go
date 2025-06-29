package middleware

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nbtca/saturday/util"
	"github.com/sirupsen/logrus"
)

// HumaLogger creates a Huma middleware for request logging
func HumaLogger() func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		startTime := time.Now()
		
		// Process request
		next(ctx)
		
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime).Microseconds()
		
		// Extract user ID from auth context if available
		var userID string
		if auth := GetAuthContext(ctx.Context()); auth != nil {
			userID = auth.ID
		}
		
		// Get response status from context
		status := ctx.Status()
		if status == 0 {
			status = 200 // Default to 200 if not set
		}
		
		util.Logger.WithFields(logrus.Fields{
			"status_code":   status,
			"latency":       latencyTime,
			"client_ip":     ctx.Header("X-Forwarded-For"), // Use forwarded IP or fallback
			"time_stamp":    endTime.Format(time.RFC1123),
			"method":        ctx.Method(),
			"path":          ctx.URL().Path,
			"request_proto": ctx.Header("X-Forwarded-Proto"),
			"user_agent":    ctx.Header("User-Agent"),
			"userId":        userID,
		}).Info("HTTP Request")
	}
}