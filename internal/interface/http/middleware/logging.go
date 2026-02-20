package middleware

import (
	"time"

	"github.com/alijayanet/gembok-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		fullPath := path
		if rawQuery != "" {
			fullPath = path + "?" + rawQuery
		}

		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", fullPath),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
		}

		// Attach gin error messages (if any handler set them)
		if errs := c.Errors.ByType(gin.ErrorTypePrivate); len(errs) > 0 {
			fields = append(fields, zap.String("errors", errs.String()))
		}

		switch {
		case statusCode >= 500:
			logger.Error("HTTP 5xx", fields...)
		case statusCode >= 400:
			logger.Warn("HTTP 4xx", fields...)
		default:
			logger.Info("HTTP 2xx/3xx", fields...)
		}
	}
}
