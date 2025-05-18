package service

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/miajio/dpsk/internal/cache"
	"github.com/miajio/dpsk/internal/dto"
	"github.com/miajio/dpsk/internal/errors"
	"go.uber.org/zap"
)

// GetLoginUser 获取登录用户信息
// 通过中间件获取登录用户信息
func GetLoginUser(ctx *gin.Context) (*dto.LoginUser, error) {
	// 获取token
	token := ctx.GetHeader("Authorization")
	if token == "" {
		return nil, errors.ErrNoToken
	}
	c := context.Background()
	loginUserStr, outTime, err := cache.RedisClient.GetWithTTL(c, cache.JWT.Prefix+token)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}
	loginUser := dto.LoginUser{}
	if err := json.Unmarshal([]byte(loginUserStr), &loginUser); err != nil {
		return nil, errors.ErrTokenParseFailed
	}

	// 校验当前时间与outTime是否接近12小时,如果接近则刷新令牌
	if cache.JWT.Expires/2 > outTime {
		if err := cache.RedisClient.SetEx(ctx, cache.JWT.Prefix+token, loginUserStr, cache.JWT.Expires); err != nil {
			zap.S().Errorf("redis续期token错误:%v", err)
			return nil, errors.ErrTokenRefreshFailed
		}
	}
	return &loginUser, nil
}
