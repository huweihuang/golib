package middlerwares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"

	log "github.com/huweihuang/golib/logger/zap"
)

func Logger() gin.HandlerFunc {
	return ZapMiddleware(log.Logger())
}

func ZapMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		httpFields, statusCode := GetHttpFields(c)

		if statusCode >= 500 {
			logger.With("httpFields", httpFields).Error()
		} else if statusCode >= 400 {
			logger.With("httpFields", httpFields).Warn()
		} else {
			logger.With("httpFields", httpFields).Info()
		}
	}
}

func LogrusMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		httpFields, statusCode := GetHttpFields(c)

		if statusCode >= 500 {
			logger.WithFields(httpFields).Error()
		} else if statusCode >= 400 {
			logger.WithFields(httpFields).Warn()
		} else {
			logger.WithFields(httpFields).Info()
		}
	}
}

func GetHttpFields(c *gin.Context) (httpFields map[string]interface{}, status int) {
	// Start timer
	start := time.Now()
	// Process request
	c.Next()
	// Stop timer
	end := time.Now()
	latency := end.Sub(start)

	statusCode := c.Writer.Status()

	httpFields = map[string]interface{}{
		"path":        c.Request.URL.Path,
		"query":       c.Request.URL.RawQuery,
		"latency":     formatLatency(latency),
		"ip":          c.ClientIP(),
		"method":      c.Request.Method,
		"status_code": statusCode,
		"req_id":      c.GetString("req_id"),
	}

	return httpFields, statusCode
}

// formatLatency convert to milliseconds
func formatLatency(latency time.Duration) int {
	return int(latency.Seconds() * 1000)
}
