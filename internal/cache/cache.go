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
	File        *FileConfig
)

type JWTConfig struct {
	Prefix  string        `toml:"prefix" yaml:"prefix"`   // token前缀
	Expires time.Duration `toml:"expires" yaml:"expires"` // token过期时间
}

type FileConfig struct {
	SavePath string `toml:"save_path" yaml:"save_path"` // 文件存储路径
	FilePath string `toml:"file_path" yaml:"file_path"` // 文件存储路径-系统文件前缀
	Url      string `toml:"url" yaml:"url"`             // 文件访问路径
}
