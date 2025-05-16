package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySqlConfig 数据库配置
type MySqlConfig struct {
	Host            string        `toml:"host" yaml:"host"`
	Port            int           `toml:"port" yaml:"port"`
	User            string        `toml:"user" yaml:"user"`
	Password        string        `toml:"password" yaml:"password"`
	Name            string        `toml:"name" yaml:"name"`
	MaxIdleConns    int           `toml:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns    int           `toml:"max_open_conns" yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `toml:"conn_max_lifetime" yaml:"conn_max_lifetime"` // 字符串解析为时间间隔
}

// getDSN 获取数据库连接字符串
func (c *MySqlConfig) getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
}

// Init 初始化数据库
func (c *MySqlConfig) Init() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(c.getDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(c.ConnMaxLifetime)
	return db, nil
}
