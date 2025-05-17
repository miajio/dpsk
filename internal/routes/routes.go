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
	Name    string
	Version string
	Port    string
	Env     string // prod/dev/test
}

type Route struct {
	cfg AppConfig
}

func NewRoute(app AppConfig) *Route {
	return &Route{app}
}

func (app *Route) SetName(name string) {
	app.cfg.Name = name
}

func (app *Route) SetVersion(version string) {
	app.cfg.Version = version
}

func (app *Route) GetPort() string {
	return app.cfg.Port
}

func (app *Route) SetupRouter(db *gorm.DB) *gin.Engine {
	if app.cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.New()

	// 中间件
	router.Use(middleware.GinLogger())
	router.Use(middleware.GinRecovery(true))
	router.Use(middleware.VersionMiddleware(app.cfg.Version))
	router.Use(middleware.ErrorHandler())
	router.Use(cors.Default())

	app.register(router, db)
	return router
}

func (app *Route) register(router *gin.Engine, db *gorm.DB) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusOK, "pong", nil, nil))
	})
	router.GET("/env", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusOK, "success", map[string]any{
			"name":    app.cfg.Name,
			"env":     app.cfg.Env,
			"version": app.cfg.Version,
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

		// 私有路由
		private := api.Group("/private").Use(middleware.TokenMiddleware())
		{
			private.GET("/logout", userController.Logout)
		}
	}
}
