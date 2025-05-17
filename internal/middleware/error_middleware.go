package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miajio/dpsk/internal/dto"
	"go.uber.org/zap"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				zap.S().Errorf("Panic: %v", err)
				c.AbortWithStatusJSON(http.StatusOK, dto.NewBaseResponse(http.StatusInternalServerError, "fail", nil, errors.New("internal server error")))
				return
			}
		}()
		c.Next()
	}
}
