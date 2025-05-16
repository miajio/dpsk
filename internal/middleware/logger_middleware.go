package middleware

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return ginzap.Ginzap(zap.L(), time.RFC3339, true)
}

// GinRecovery recover掉请求中的panic，并使用zap记录
func GinRecovery(stack bool) gin.HandlerFunc {
	return ginzap.RecoveryWithZap(zap.L(), stack)
}
