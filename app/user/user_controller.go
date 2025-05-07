package user

import (
	"github.com/gin-gonic/gin"
	"github.com/miajio/dpsk/pkg/router"
)

type UserControllerTemp interface {
	Register(ctx *gin.Context)
}

type UserControllerImpl struct {
}

var UserController router.Controller = &UserControllerImpl{}

// Option
func (c *UserControllerImpl) Option(e *gin.Engine) {
	group := e.Group("/user")
	{
		group.POST("/register", c.Register)
	}
}

func (c *UserControllerImpl) Register(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "register",
	})
}
