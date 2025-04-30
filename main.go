package main

import (
	"github.com/gin-gonic/gin"
	"github.com/miajio/dpsk/comm/ctx"
	"github.com/miajio/dpsk/comm/log"
	"github.com/miajio/dpsk/middleware"
	"github.com/miajio/dpsk/pkg/config"
)

var cfg = Config{}

func init() {
	if err := config.ReadTomlConfig("./config.toml", &cfg); err != nil {
		panic(err)
	}

	if logger, err := cfg.Log.Generator(); err != nil {
		panic(err)
	} else {
		log.SetLogger(logger)
	}

	if redis, err := cfg.Redis.Generator(); err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	} else {
		ctx.Redis = redis
	}

	if db, err := cfg.Database.Generator(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	} else {
		ctx.DB = db
	}

	// 设置Gin模式
	if cfg.App.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

func main() {
	// 创建Gin路由
	r := gin.New()

	// 启动服务器
	port := cfg.App.Port
	if port == "" {
		port = ":8080"
	}
	log.Infof("服务启动，监听端口 %s", port)

	r.Use(
		middleware.CORS(),
		middleware.Recovery(),
	)

	if err := r.Run(port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
