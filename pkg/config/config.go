package config

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ReadConfig 读取配置文件
func ReadTomlConfig(configPath string, rawVal any) error {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("获取配置文件绝对路径失败: %w", err)
	}

	viper.SetConfigFile(absPath)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	if err := viper.Unmarshal(rawVal); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 监听配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("配置文件变更，重新加载:", e.Name)
		if err := viper.Unmarshal(rawVal); err != nil {
			log.Printf("重新加载配置文件失败: %v\n", err)
		}
	})

	return nil
}
