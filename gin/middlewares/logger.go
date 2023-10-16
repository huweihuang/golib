package middlerwares

import (
	"time"

	log "github.com/huweihuang/golib/logger/logrus"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	return LogMiddleware(log.Log())
}

func LogMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		statusCode := c.Writer.Status()

		fields := logrus.Fields{
			"path":    path,
			"query":   query,
			"latency": formatLatency(latency),
			"ip":      c.ClientIP(),
			"method":  c.Request.Method,
			"code":    statusCode,
			"req_id":  c.GetString("req_id"),
		}

		if statusCode >= 500 {
			logger.WithFields(fields).Error()
		} else if statusCode >= 400 {
			logger.WithFields(fields).Warn()
		} else {
			logger.WithFields(fields).Info()
		}
	}
}

// formatLatency convert to milliseconds
func formatLatency(latency time.Duration) int {
	return int(latency.Seconds() * 1000)
}
