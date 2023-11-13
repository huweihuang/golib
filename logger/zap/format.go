package zap

import (
	"fmt"

	"go.uber.org/zap"
)

// Fields type, used to pass to `With`.
type Fields map[string]interface{}

func With(k string, v interface{}) zap.Field {
	return zap.Any(k, v)
}

func WithError(err error) zap.Field {
	return zap.NamedError("error", err)
}

func Info(msg string, args ...zap.Field) {
	L.Info(msg, args...)
}

func Error(msg string, args ...zap.Field) {
	L.Error(msg, args...)
}

func Warn(msg string, args ...zap.Field) {
	L.Warn(msg, args...)
}

func Debug(msg string, args ...zap.Field) {
	L.Debug(msg, args...)
}

func Fatal(msg string, args ...zap.Field) {
	L.Fatal(msg, args...)
}

func Infof(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	L.Info(logMsg)
}

func Errorf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	L.Error(logMsg)
}

func Warnf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	L.Warn(logMsg)
}

func Debugf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	L.Debug(logMsg)
}

func Fatalf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	L.Fatal(logMsg)
}
