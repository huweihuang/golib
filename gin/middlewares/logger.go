package middlerwares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"

	log "github.com/huweihuang/golib/logger/zap"
)

func Logger() gin.HandlerFunc {
	return ZapMiddleware(log.Log())
}

func ZapMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		// Process request
		c.Next()
		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		statusCode := c.Writer.Status()

		httpFields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.Int("latency", formatLatency(latency)),
			zap.String("req_id", c.GetString("req_id")),
			zap.String("user-agent", c.Request.UserAgent()),
		}

		if statusCode >= 500 {
			logger.With(httpFields...).Error("http fields")
		} else if statusCode >= 400 {
			logger.With(httpFields...).Warn("http fields")
		} else {
			logger.With(httpFields...).Info("http fields")
		}
	}
}

func LogrusMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		// Process request
		c.Next()
		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		statusCode := c.Writer.Status()

		httpFields := logrus.Fields{
			"path":        c.Request.URL.Path,
			"query":       c.Request.URL.RawQuery,
			"latency":     formatLatency(latency),
			"ip":          c.ClientIP(),
			"method":      c.Request.Method,
			"status_code": statusCode,
			"req_id":      c.GetString("req_id"),
			"user-agent":  c.Request.UserAgent(),
		}

		if statusCode >= 500 {
			logger.WithFields(httpFields).Error()
		} else if statusCode >= 400 {
			logger.WithFields(httpFields).Warn()
		} else {
			logger.WithFields(httpFields).Info()
		}
	}
}

// formatLatency convert to milliseconds
func formatLatency(latency time.Duration) int {
	return int(latency.Seconds() * 1000)
}
