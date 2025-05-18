package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miajio/dpsk/internal/dto"
	"github.com/miajio/dpsk/internal/service"
	"go.uber.org/zap"
)

func TokenMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if loginUser, err := service.GetLoginUser(ctx); err != nil {
			zap.S().Errorf("获取登录用户信息失败:%v", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.NewBaseResponse(http.StatusUnauthorized, "未授权", nil, err))
			return
		} else {
			ctx.Set("loginUser", loginUser)
		}
		ctx.Next()
	}
}
