package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger is an alias to gin.HandlerFunc
type Logger gin.HandlerFunc

// NewLogger creates a logger
func NewLogger(logger *zap.Logger) Logger {
	logger.Info("Creating logger")

	namedLogger := logger.Named("logger")

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		timeStamp := time.Now()

		namedLogger.Info("Request",
			zap.String("path", c.Request.URL.Path),
			zap.String("raw_query", c.Request.URL.RawQuery),
			zap.String("full_path", c.FullPath()),
			zap.Duration("latency", timeStamp.Sub(start)),
			zap.String("client_ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.Int("status", c.Writer.Status()),
			zap.String("error_message", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Int("body_size", c.Writer.Size()),
		)
	}
}
