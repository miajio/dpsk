package cache

import (
	"github.com/miajio/dpsk/internal/routes"
	"github.com/miajio/dpsk/pkg/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	RedisClient *redis.RedisClient
	DB          *gorm.DB
	Log         *zap.Logger
	AppConfig   *routes.AppConfig
)
