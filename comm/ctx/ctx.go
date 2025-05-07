package ctx

import (
	"github.com/miajio/dpsk/pkg/cache"
	"gorm.io/gorm"
)

var (
	File  *FileConfig
	DB    *gorm.DB
	Redis *cache.RedisClient
)

type FileConfig struct {
	UploadDir         string `yaml:"uploadDir"`         // 上传文件保存路径
	BaseUrl           string `yaml:"baseUrl"`           // 上传文件访问路径
	MaxUploadSize     int64  `yaml:"maxUploadSize"`     // 最大上传文件大小
	MinChunkSize      int64  `yaml:"minChunkSize"`      // 分片大小
	StreamDownloadMin int64  `yaml:"streamDownloadMin"` // 启用流式下载的最小文件大小
}
