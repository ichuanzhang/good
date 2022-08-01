package middleware

import (
	"github.com/gin-gonic/gin"
	"good/internal/core"
	"good/pkg/limiter"
	"net/http"
	"sync"
)

// tokenBuffer 令牌缓冲区，只在分布式限流模式下才使用到，每次批量获取令牌存入本地缓冲区
var tokenBuffer = &TokenBuffer{}

type TokenBuffer struct {
	count int
	mu    sync.Mutex
}

// Limiter 限流中间件
func Limiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ok bool
		switch core.Limiter.GetMode() {
		case limiter.ModeSingle:
			ok = singleLimiter()
		case limiter.ModeDistributed:
			ok = distributedLimiter()
		default:
			panic("not support limiter mode")
		}

		if !ok {
			c.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}

		c.Next()
		return
	}
}

// singleLimiter 单节点限流
func singleLimiter() bool {
	_, ok := core.Limiter.GetToken(1)
	return ok
}

// singleLimiter 分布式限流
func distributedLimiter() bool {
	tokenBuffer.mu.Lock()
	defer tokenBuffer.mu.Unlock()

	if tokenBuffer.count > 0 {
		tokenBuffer.count--
		return true
	}

	if n, ok := core.Limiter.GetToken(core.Limiter.GetRate()); ok {
		tokenBuffer.count += n - 1
		return true
	}
	return false
}
