package log

import (
	"fmt"

	"go.uber.org/zap"
)

var log *zap.Logger

func SetLogger(logger *zap.Logger) {
	log = logger
}

func GetLogger() *zap.Logger {
	return log
}

func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

func Debugf(template string, args ...interface{}) {
	log.Debug(fmt.Sprintf(template, args...))
}

func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

func Fatalf(template string, args ...interface{}) {
	log.Fatal(fmt.Sprintf(template, args...))
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Errorf(template string, args ...interface{}) {
	log.Error(fmt.Sprintf(template, args...))
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

func Warnf(template string, args ...interface{}) {
	log.Warn(fmt.Sprintf(template, args...))
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Infof(template string, args ...interface{}) {
	log.Info(fmt.Sprintf(template, args...))
}
