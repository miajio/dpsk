package config

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ReloadFunc 定义配置文件重载时的回调函数类型
type ReloadFunc func(config any) error

// ConfigManager 配置管理器
type ConfigManager struct {
	viper      *viper.Viper
	configPath string
	readCfgObj any
	reloadFunc ReloadFunc
	mu         sync.Mutex
}

// NewConfigManager 创建新的配置管理器
func NewConfigManager(configPath string, readCfgObj any, reloadFunc ReloadFunc) (*ConfigManager, error) {
	// 获取绝对路径
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %v", err)
	}

	// 初始化 viper
	v := viper.New()
	v.SetConfigFile(absPath)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}

	// 解析配置到目标对象
	if err := v.Unmarshal(readCfgObj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// 执行初始回调
	if reloadFunc != nil {
		if err := reloadFunc(readCfgObj); err != nil {
			return nil, fmt.Errorf("initial reload func failed: %v", err)
		}
	}

	cm := &ConfigManager{
		viper:      v,
		configPath: absPath,
		readCfgObj: readCfgObj,
		reloadFunc: reloadFunc,
	}

	// 设置配置文件变化监听
	v.OnConfigChange(func(e fsnotify.Event) {
		cm.onConfigChange()
	})
	v.WatchConfig()

	return cm, nil
}

// onConfigChange 处理配置文件变化
func (cm *ConfigManager) onConfigChange() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	log.Println("Config file changed, reloading...")

	// 重新读取配置
	if err := cm.viper.ReadInConfig(); err != nil {
		log.Printf("Failed to read config: %v\n", err)
		return
	}

	// 解析到配置对象
	if err := cm.viper.Unmarshal(cm.readCfgObj); err != nil {
		log.Printf("Failed to unmarshal config: %v\n", err)
		return
	}

	// 执行回调函数
	if cm.reloadFunc != nil {
		if err := cm.reloadFunc(cm.readCfgObj); err != nil {
			log.Printf("Reload func failed: %v\n", err)
		}
	}
}

// GetConfig 获取当前配置对象
func (cm *ConfigManager) GetConfig() any {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.readCfgObj
}
