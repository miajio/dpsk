package main

import (
	"github.com/miajio/dpsk/pkg/cache"
	"github.com/miajio/dpsk/pkg/database"
	"github.com/miajio/dpsk/pkg/logger"
)

type Config struct {
	App      AppConfig               `toml:"app" json:"app"`           // 应用配置
	Log      logger.LogConfig        `toml:"log" json:"log"`           // 日志配置
	SysFile  SysFileConfig           `toml:"sysFile" json:"sysFile"`   // 文件配置
	Redis    cache.RedisConfig       `toml:"redis" json:"redis"`       // redis配置
	Database database.DatabaseConfig `toml:"database" json:"database"` // 数据库配置
}

type AppConfig struct {
	Name    string `toml:"name" json:"name"`       // 应用名称
	Version string `toml:"version" json:"version"` // 应用版本
	Port    string `toml:"port" json:"port"`       // 应用端口
	Env     string `toml:"env" json:"env"`         // 应用环境
}

type SysFileConfig struct {
	UploadDir         string `yaml:"uploadDir"`         // 上传文件保存路径
	BaseUrl           string `yaml:"baseUrl"`           // 上传文件访问路径
	MaxUploadSize     int64  `yaml:"maxUploadSize"`     // 最大上传文件大小
	MinChunkSize      int64  `yaml:"minChunkSize"`      // 分片大小
	StreamDownloadMin int64  `yaml:"streamDownloadMin"` // 启用流式下载的最小文件大小
}
