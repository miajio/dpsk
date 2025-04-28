package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig 日志配置
type LogConfig struct {
	Level      string `toml:"level" json:"level"`           // 日志级别 debug/info/warn/error/fatal
	Format     string `toml:"format" json:"format"`         // 日志格式 json/text
	Output     string `toml:"output" json:"output"`         // 输出方式 file/stdout/both
	Path       string `toml:"path" json:"path"`             // 日志路径
	MaxSize    int    `toml:"maxSize" json:"maxSize"`       // 日志文件大小 MB
	MaxBackups int    `toml:"maxBackups" json:"maxBackups"` // 日志文件数量
	MaxAge     int    `toml:"maxAge" json:"maxAge"`         // 日志文件过期时间 days
	Compress   bool   `toml:"compress" json:"compress"`     // 是否压缩
}

// Generator 生成日志实例
func (logConfig *LogConfig) Generator() (*zap.Logger, error) {
	var core zapcore.Core

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置日志级别
	var level zapcore.Level
	switch logConfig.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	// 输出方式
	var outputs []zapcore.WriteSyncer
	switch logConfig.Output {
	case "file":
		// 确保日志目录存在
		if err := os.MkdirAll(filepath.Dir(logConfig.Path), os.ModePerm); err != nil {
			return nil, err
		}

		lumberJackLogger := &lumberjack.Logger{
			Filename:   logConfig.Path,
			MaxSize:    logConfig.MaxSize,
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge,
			Compress:   logConfig.Compress,
		}
		outputs = append(outputs, zapcore.AddSync(lumberJackLogger))
	case "stdout":
		outputs = append(outputs, zapcore.AddSync(os.Stdout))
	case "both":
		if err := os.MkdirAll(filepath.Dir(logConfig.Path), os.ModePerm); err != nil {
			return nil, err
		}

		lumberJackLogger := &lumberjack.Logger{
			Filename:   logConfig.Path,
			MaxSize:    logConfig.MaxSize,
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge,
			Compress:   logConfig.Compress,
		}
		outputs = append(outputs, zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout))
	}

	// 创建核心
	if logConfig.Format == "json" {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(outputs...),
			level,
		)
	} else {
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(outputs...),
			level,
		)
	}

	// 创建Logger
	log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// 替换全局Logger
	zap.ReplaceGlobals(log)
	return log, nil
}
