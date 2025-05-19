package routes

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/miajio/dpsk/internal/controller"
	"github.com/miajio/dpsk/internal/dto"
	"github.com/miajio/dpsk/internal/middleware"
	"github.com/miajio/dpsk/internal/repository"
	"github.com/miajio/dpsk/internal/service"
	"go.uber.org/zap"
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
	cfg       AppConfig    // config
	server    *http.Server // http server
	engine    *gin.Engine  // gin engine
	db        *gorm.DB     // 数据库连接
	isReady   atomic.Value // 是否就绪
	isHealthy atomic.Value // 是否健康
}

func NewRoute(app AppConfig, db *gorm.DB) *Route {
	route := &Route{cfg: app, db: db}
	route.isReady.Store(false)
	route.isHealthy.Store(true)
	return route
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

func (app *Route) SetupRouter() {
	if app.cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	app.engine = gin.New()

	server := &http.Server{
		Addr:    app.cfg.Port,
		Handler: app.engine,
	}

	app.server = server

	// 中间件
	app.engine.Use(middleware.GinLogger())
	app.engine.Use(middleware.GinRecovery(true))
	app.engine.Use(middleware.VersionMiddleware(app.cfg.Version))
	app.engine.Use(middleware.ErrorHandler())
	app.engine.Use(cors.Default())

	app.register()
}

func (app *Route) register() {
	app.engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusOK, "pong", nil, nil))
	})
	app.engine.GET("/env", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusOK, "success", map[string]any{
			"name":    app.cfg.Name,
			"env":     app.cfg.Env,
			"version": app.cfg.Version,
		}, nil))
	})

	// 健康检查路由
	app.engine.GET("/ready", func(ctx *gin.Context) {
		if app.isReady.Load().(bool) {
			ctx.JSON(http.StatusOK, gin.H{"status": "ready"})
		} else {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
		}
	})

	// 存活探针
	app.engine.GET("/health", func(c *gin.Context) {
		if app.isHealthy.Load().(bool) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
		}
	})

	userController := controller.NewUserController(service.NewUserService(repository.NewUserRepository(app.db)))
	fileController := controller.NewFileController(service.NewFileService(repository.NewFileRepository(app.db)))

	// API路由组
	api := app.engine.Group("/api")
	{
		// 公共路由
		authPublic := api.Group("/auth")
		{
			authPublic.POST("/login", userController.Login)
			authPublic.POST("/register", userController.Register)
		}

		// 文件路由
		filePublic := api.Group("/file")
		{
			filePublic.POST("/upload", fileController.Upload)
		}

		// 私有路由
		private := api.Group("/private").Use(middleware.TokenMiddleware())
		{
			private.GET("/logout", userController.Logout)
		}
	}
}

func (app *Route) Start() {
	app.isReady.Store(true)
	// 启动服务器
	go func() {
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.S().Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 优雅关闭处理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("关闭服务中...")
	// 标记服务为不健康
	app.isHealthy.Store(false)
	app.isReady.Store(false)

	// 等待一段时间让负载均衡器检测到状态变化
	time.Sleep(5 * time.Second)
	// 关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.server.Shutdown(ctx); err != nil {
		zap.S().Fatalf("服务关闭失败:%v", err)
	}

	zap.L().Info("服务已关闭")
}
