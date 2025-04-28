package router

import "github.com/gin-gonic/gin"

// Controller 路由控制器
type Controller interface {
	Option(e *gin.Engine)
}

// Router 路由管理器
type Router struct {
	controllers []Controller
}

// NewRouter 创建路由管理器
func NewRouter() *Router {
	return &Router{
		controllers: make([]Controller, 0),
	}
}

// AddController 添加控制器
func (r *Router) AddController(c ...Controller) {
	r.controllers = append(r.controllers, c...)
}

// RegisterAllController 路由注册所有controller
func (r *Router) RegisterAllController(e *gin.Engine) {
	for _, c := range r.controllers {
		c.Option(e)
	}
}
