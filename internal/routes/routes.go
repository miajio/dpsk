package routes

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/miajio/dpsk/internal/controller"
	"github.com/miajio/dpsk/internal/dto"
	"github.com/miajio/dpsk/internal/middleware"
	"github.com/miajio/dpsk/internal/repository"
	"github.com/miajio/dpsk/internal/service"
	"gorm.io/gorm"
)

// App 应用基础配置
type AppConfig struct {
	Name    string `toml:"name" yaml:"name"`
	Version string `toml:"version" yaml:"version"`
	Port    string `toml:"port" yaml:"port"`
	Env     string `toml:"env" yaml:"env"` // prod/dev/test
}

func (app *AppConfig) SetupRouter(db *gorm.DB) *gin.Engine {
	if app.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.New()

	// 中间件
	router.Use(middleware.VersionMiddleware(app.Version))
	router.Use(middleware.ErrorHandler())
	router.Use(cors.Default())
	app.register(router, db)
	return router
}

func (app *AppConfig) register(router *gin.Engine, db *gorm.DB) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusOK, "pong", nil, nil))
	})
	router.GET("/env", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusOK, "success", map[string]any{
			"name":    app.Name,
			"env":     app.Env,
			"version": app.Version,
		}, nil))
	})

	userController := controller.NewUserController(service.NewUserService(repository.NewUserRepository(db)))

	// API路由组
	api := router.Group("/api")
	{
		// 公共路由
		public := api.Group("/auth")
		{
			public.POST("/login", userController.Login)
			public.POST("/register", userController.Register)
		}
	}
}
