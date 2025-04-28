package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `toml:"host" json:"host"`                       // 数据库地址
	Port            int           `toml:"port" json:"port"`                       // 数据库端口
	User            string        `toml:"user" json:"user"`                       // 数据库用户名
	Password        string        `toml:"password" json:"password"`               // 数据库密码
	Name            string        `toml:"name" json:"name"`                       // 数据库名称
	MaxIdleConns    int           `toml:"maxIdleConns" json:"maxIdleConns"`       // 最大空闲连接数 10
	MaxOpenConns    int           `toml:"maxOpenConns" json:"maxOpenConns"`       // 最大打开连接数 100
	ConnMaxLifetime time.Duration `toml:"connMaxLifetime" json:"connMaxLifetime"` // 连接最大生命周期 5m
}

// GetDSN 获取dsn
func (dbConfig *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)
}

// Generator 生成gorm连接
func (dbConfig *DatabaseConfig) Generator() (*gorm.DB, error) {
	var err error
	db, err := gorm.Open(mysql.Open(dbConfig.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)

	return db, nil
}
