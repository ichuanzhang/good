package router

import (
	"github.com/gin-gonic/gin"
	"good/internal/config"
	"good/internal/controller"
	"good/internal/router/middleware"
)

var middlewareMap = map[string]gin.HandlerFunc{
	"recovery": gin.Recovery(),
	"cross":    middleware.Cors(),
	"traceLog": middleware.TraceLog(),
	"limiter":  middleware.Limiter(),
}

// Init 初始化router
func Init() *gin.Engine {
	setMode()
	r := gin.New()
	initMiddleware(r)
	initController(r)
	return r
}

// setMode 设置gin启动模式
func setMode() {
	switch config.C.Application.Mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	}
}

// initMiddleware 注册Middleware
func initMiddleware(router *gin.Engine) {
	for _, v := range config.C.Middleware {
		f, ok := middlewareMap[v]
		if !ok {
			continue
		}
		router.Use(f)
	}
}

// initController 注册API
func initController(router *gin.Engine) {
	system := router.Group("")
	{
		system.GET("/ping", controller.Ping)
	}
}
