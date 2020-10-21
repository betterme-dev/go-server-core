package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Logger is a middleware which logs all HTTP requests. Slightly modified version of gin.Logger.
func Logger(c *gin.Context) {
	// Start timer
	start := time.Now()
	path := c.Request.URL.String()
	userAgent := c.Request.Header.Get("user-agent")

	// Process request
	c.Next()

	// Stop timer
	latency := time.Since(start)

	method := c.Request.Method
	statusCode := c.Writer.Status()
	errors := c.Errors.String()
	if errors != "" {
		errors = errors[:len(errors)-1]
	}

	logger := log.WithFields(log.Fields{
		"status":  statusCode,
		"latency": latency.Seconds(),
		"remote":  c.ClientIP(),
		"ua":      userAgent,
		"method":  method,
		"size":    c.Writer.Size(),
	})

	if errors != "" {
		logger.WithFields(log.Fields{
			"error": errors,
		}).Error(path)
	} else {
		// all codes below 400 and 404 not critical errors so we log them with error level "info"
		if statusCode < http.StatusBadRequest || statusCode == http.StatusNotFound {
			logger.Info(path)
		} else {
			logger.Error(path)
		}
	}
}
