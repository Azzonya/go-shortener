package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

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
		)

		logger.Info("Response",
			zap.Int("statusCode", statusCode),
			zap.Int64("contentLength", contentLength),
		)
	}
}
