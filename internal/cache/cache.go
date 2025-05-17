package cache

import (
	"time"

	"github.com/miajio/dpsk/pkg/redis"
	"gorm.io/gorm"
)

var (
	RedisClient *redis.RedisClient
	DB          *gorm.DB
	JWT         *JWTConfig
)

type JWTConfig struct {
	Prefix  string        `toml:"prefix" yaml:"prefix"`   // token前缀
	Expires time.Duration `toml:"expires" yaml:"expires"` // token过期时间
}
