package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `toml:"level" yaml:"level"`             // debug, info, warn, error, dpanic, panic, fatal
	Path       string `toml:"path" yaml:"path"`               // 日志文件路径
	Console    bool   `toml:"console" yaml:"console"`         // 是否输出到控制台
	File       bool   `toml:"file" yaml:"file"`               // 是否输出到文件
	MaxSize    int    `toml:"max_size" yaml:"max_size"`       // 单个日志文件最大大小(MB)
	MaxBackups int    `toml:"max_backups" yaml:"max_backups"` // 保留的旧日志文件最大数量
	MaxAge     int    `toml:"max_age" yaml:"max_age"`         // 保留旧日志文件的最大天数
	Compress   bool   `toml:"compress" yaml:"compress"`       // 是否压缩/归档旧日志文件
}

func (cfg *LoggerConfig) Init() *zap.Logger {
	// 设置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	// 自定义时间编码格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(level.CapitalString())
	}

	// 自定义调用者显示
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(caller.TrimmedPath())
	}

	// 创建Encoder配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    customLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   customCallerEncoder,
	}

	// 设置日志输出
	var cores []zapcore.Core

	// 控制台输出
	if cfg.Console {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))
	}

	// 文件输出
	if cfg.File {
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Path + "/app.log",
			MaxSize:    cfg.MaxSize,    // 日志文件最大大小(MB)
			MaxBackups: cfg.MaxBackups, // 保留旧文件的最大个数
			MaxAge:     cfg.MaxAge,     // 保留旧文件的最大天数
			Compress:   cfg.Compress,   // 是否压缩/归档旧文件
		})
		cores = append(cores, zapcore.NewCore(fileEncoder, fileWriter, level))
	}

	// 创建核心
	core := zapcore.NewTee(cores...)

	// 创建Logger
	log := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	// 替换zap全局Logger
	zap.ReplaceGlobals(log)
	return log
}
