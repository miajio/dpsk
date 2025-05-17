package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miajio/dpsk/internal/cache"
	"github.com/miajio/dpsk/internal/dto"
	"go.uber.org/zap"
)

func TokenMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取token
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusOK, dto.NewBaseResponse(http.StatusUnauthorized, "未提供认证令牌", nil, nil))
			return
		}
		c := context.Background()
		loginUserStr, outTime, err := cache.RedisClient.GetWithTTL(c, cache.JWT.Prefix+token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, dto.NewBaseResponse(http.StatusUnauthorized, "认证令牌无效", nil, nil))
			return
		}
		loginUser := dto.LoginUser{}
		if err := json.Unmarshal([]byte(loginUserStr), &loginUser); err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, dto.NewBaseResponse(http.StatusUnauthorized, "令牌解析失败", nil, nil))
			return
		}

		// 校验当前时间与outTime是否接近12小时,如果接近则刷新令牌
		if cache.JWT.Expires/2 > outTime {
			if err := cache.RedisClient.SetEx(ctx, cache.JWT.Prefix+token, loginUserStr, cache.JWT.Expires); err != nil {
				zap.S().Errorf("redis续期token错误:%v", err)
				ctx.AbortWithStatusJSON(http.StatusOK, dto.NewBaseResponse(http.StatusUnauthorized, "令牌续期失败", nil, err))
				return
			}
		}

		ctx.Set("loginUser", loginUser)
		ctx.Next()
	}
}
