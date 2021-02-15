package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const loggerKey = "logger"

// Middleware return middleware function
func Middleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		l := logger.With((zap.String("id", c.GetHeader("X-Request-ID"))))
		c.Set(loggerKey, l)

		c.Next()

		latency := time.Since(t)
		l.Info("request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
		)
	}
}

// Extract extract zap logger from context
func Extract(c *gin.Context) *zap.Logger {
	l, ok := c.Get(loggerKey)
	if ok {
		return l.(*zap.Logger)
	}
	return zap.NewExample()
}
