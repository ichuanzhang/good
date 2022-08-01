package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"good/internal/core"
	"net/http"
	"time"
)

// TraceLog 日志中间件
func TraceLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		if c.Writer.Status() == http.StatusOK {
			core.Logger.Info("OK",
				zap.String("method", c.Request.Method),
				zap.Int("status", c.Writer.Status()),
				zap.String("url", c.Request.URL.String()),
				zap.Duration("duration", time.Since(start)),
			)
		} else {
			core.Logger.Error("ERROR",
				zap.String("method", c.Request.Method),
				zap.Int("status", c.Writer.Status()),
				zap.String("url", c.Request.URL.String()),
				zap.Duration("duration", time.Since(start)),
			)
		}
	}
}
