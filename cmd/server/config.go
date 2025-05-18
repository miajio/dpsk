package main

import (
	"github.com/miajio/dpsk/internal/cache"
	"github.com/miajio/dpsk/internal/routes"
	"github.com/miajio/dpsk/pkg/database"
	"github.com/miajio/dpsk/pkg/logger"
	"github.com/miajio/dpsk/pkg/redis"
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
