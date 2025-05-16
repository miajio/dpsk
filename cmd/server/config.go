package main

import (
	"github.com/miajio/dpsk/internal/routes"
	"github.com/miajio/dpsk/pkg/database"
	"github.com/miajio/dpsk/pkg/logger"
	"github.com/miajio/dpsk/pkg/redis"
)

// AppConfig 应用整体配置结构
type AppConfig struct {
	App      routes.AppConfig     `toml:"app"`
	Log      logger.LoggerConfig  `toml:"log"`
	Database database.MySqlConfig `toml:"database"`
	Redis    redis.RedisConfig    `toml:"redis"`
}
