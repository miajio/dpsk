package middleware

import "github.com/gin-gonic/gin"

// VersionMiddleware 版本中间件
func VersionMiddleware(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("api-version", version)
		c.Next()
	}
}
