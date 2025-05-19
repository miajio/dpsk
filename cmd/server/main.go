package main

import (
	"fmt"
	"log"
	"os"

	"github.com/miajio/dpsk/internal/cache"
	"github.com/miajio/dpsk/internal/model"
	"github.com/miajio/dpsk/internal/routes"
	"github.com/miajio/dpsk/pkg/config"
	"github.com/miajio/dpsk/pkg/database"
	"github.com/miajio/dpsk/pkg/logger"
	"github.com/miajio/dpsk/pkg/redis"
	"go.uber.org/zap"
)

// Config 应用整体配置结构
type Config struct {
	App      routes.AppConfig     `toml:"app" yaml:"app"`
	Log      logger.LoggerConfig  `toml:"log" yaml:"log"`
	Database database.MySqlConfig `toml:"database" yaml:"database"`
	Redis    redis.RedisConfig    `toml:"redis" yaml:"redis"`
	JWT      cache.JWTConfig      `toml:"jwt" yaml:"jwt"`
	File     cache.FileConfig     `toml:"file" yaml:"file"`
}

var (
	isInit bool
	route  *routes.Route
)

func cfgInit(config any) error {
	c, ok := config.(*Config)
	if !ok {
		return fmt.Errorf("配置类型无效")
	}
	if !isInit {
		c.Log.Init()
		var err error
		cache.DB, err = c.Database.Init()
		if err != nil {
			zap.S().Fatalf("数据库初始化失败: %v", err)
		} else {
			zap.L().Info("数据库初始化成功")
		}
		if err := cache.DB.AutoMigrate(model.UserModel{}, model.FileModel{}); err != nil {
			zap.S().Fatalf("数据库迁移失败: %v", err)
		} else {
			zap.L().Info("数据库迁移成功")
		}
		cache.RedisClient, err = c.Redis.Init()
		if err != nil {
			zap.S().Fatalf("Redis初始化失败: %v", err)
		} else {
			zap.L().Info("Redis初始化成功")
		}

		route = routes.NewRoute(routes.AppConfig{
			Name:    c.App.Name,
			Version: c.App.Version,
			Port:    c.App.Port,
			Env:     c.App.Env,
		}, cache.DB)

		cache.JWT = &c.JWT
		cache.File = &c.File

		// 文件目录创建
		if err := os.MkdirAll(c.File.SavePath, os.ModePerm); err != nil {
			zap.S().Fatalf("文件目录创建失败: %v", err)
		}

		isInit = true
	}

	route.SetName(c.App.Name)
	route.SetVersion(c.App.Version)
	return nil
}

func init() {
	cfg := Config{}
	// 获取运行目录
	_, err := config.NewConfigManager("config.toml", &cfg, cfgInit)
	if err != nil {
		log.Fatalf("Failed to initialize config manager: %v", err)
	}
}

func main() {
	route.SetupRouter()
	route.Start()
}
