// Package middleware provides HTTP middleware for the URL shortener application.
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger returns a Gin middleware function that logs incoming HTTP requests
// and outgoing HTTP responses using the provided logger.
func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		contentLength := int64(c.Writer.Size())

		logger.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("path1", c.Request.URL.Path),
			zap.Duration("duration2", duration),
			zap.Any("cookie", c.Request.Cookies()),
		)

		logger.Info("Response",
			zap.Int("statusCode", statusCode),
			zap.Int64("contentLength", contentLength),
		)
	}
}
