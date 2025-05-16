package main

import (
	"fmt"
	"log"

	"github.com/miajio/dpsk/internal/cache"
	"github.com/miajio/dpsk/internal/model"
	"github.com/miajio/dpsk/pkg/config"
	"go.uber.org/zap"
)

func cfgInit(config any) error {
	c, ok := config.(*AppConfig)
	if !ok {
		return fmt.Errorf("配置类型无效")
	}
	cache.Log = c.Log.Init()
	var err error
	cache.DB, err = c.Database.Init()
	if err != nil {
		zap.L().Sugar().Fatalf("数据库初始化失败: %v", err)
		// log.Fatalf("数据库初始化失败: %v", err)
	}
	if err := cache.DB.AutoMigrate(model.UserModel{}); err != nil {
		zap.L().Sugar().Fatalf("数据库迁移失败: %v", err)
		// log.Fatalf("数据库迁移失败: %v", err)
	}
	cache.RedisClient, err = c.Redis.Init()
	if err != nil {
		zap.L().Sugar().Fatalf("Redis初始化失败: %v", err)
		// log.Fatalf("Redis初始化失败: %v", err)
	}

	cache.AppConfig = &c.App
	return nil
}

func init() {
	cfg := AppConfig{}
	cm, err := config.NewConfigManager("../../config.toml", &cfg, cfgInit)
	if err != nil {
		log.Fatalf("Failed to initialize config manager: %v", err)
	}
	cfgInit(cm.GetConfig())
}

func main() {
	route := cache.AppConfig.SetupRouter(cache.DB)
	// 启动服务器
	if err := route.Run(cache.AppConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
