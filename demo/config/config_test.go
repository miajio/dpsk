package config_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/miajio/dpsk/pkg/config"
)

// AppConfig 应用整体配置结构
type AppConfig struct {
	App      App      `toml:"app"`
	Log      Log      `toml:"log"`
	Database Database `toml:"database"`
	Redis    Redis    `toml:"redis"`
}

// App 应用基础配置
type App struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
	Port    string `toml:"port"`
	Env     string `toml:"env"` // prod/dev/test
}

// Log 日志配置
type Log struct {
	Level      string `toml:"level"`  // debug/info/warn/error/fatal
	Format     string `toml:"format"` // json/text
	Output     string `toml:"output"` // file/stdout/both
	Path       string `toml:"path"`
	MaxSize    int    `toml:"max_size"`    // MB
	MaxBackups int    `toml:"max_backups"` // 数量
	MaxAge     int    `toml:"max_age"`     // 天数
	Compress   bool   `toml:"compress"`
}

// Database 数据库配置
type Database struct {
	Host            string        `toml:"host"`
	Port            int           `toml:"port"`
	User            string        `toml:"user"`
	Password        string        `toml:"password"`
	Name            string        `toml:"name"`
	MaxIdleConns    int           `toml:"max_idle_conns"`
	MaxOpenConns    int           `toml:"max_open_conns"`
	ConnMaxLifetime time.Duration `toml:"conn_max_lifetime"` // 字符串解析为时间间隔
}

// Redis 配置
type Redis struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
	PoolSize int    `toml:"pool_size"`
}

func TestRead(t *testing.T) {
	var cfg AppConfig

	reload := func(config any) error {
		c, ok := config.(*AppConfig)
		if !ok {
			return fmt.Errorf("invalid config type")
		}
		log.Printf("Config reloaded! New values: %+v", c)
		// 这里可以添加你的业务逻辑，比如重新初始化数据库连接等
		return nil
	}

	// 创建配置管理器
	cm, err := config.NewConfigManager("config.toml", &cfg, reload)
	if err != nil {
		log.Fatalf("Failed to initialize config manager: %v", err)
	}

	// 模拟长时间运行的应用
	for {
		currentCfg := cm.GetConfig().(*AppConfig)
		fmt.Printf("Current config: %+v\n", currentCfg)
		time.Sleep(10 * time.Second)
	}

}
