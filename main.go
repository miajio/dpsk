package main

import (
	"github.com/miajio/dpsk/ctx"
	"github.com/miajio/dpsk/log"
	"github.com/miajio/dpsk/pkg/config"
)

func main() {
	cfg := Config{}
	if err := config.ReadTomlConfig("./config.toml", &cfg); err != nil {
		panic(err)
	}

	if logger, err := cfg.Log.Generator(); err != nil {
		panic(err)
	} else {
		log.SetLogger(logger)
	}

	if db, err := cfg.Database.Generator(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	} else {
		ctx.DB = db
	}

}
