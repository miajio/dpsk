package ctx

import (
	"github.com/miajio/dpsk/pkg/cache"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	Redis *cache.RedisClient
)
