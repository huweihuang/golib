package zap

import (
	"fmt"

	"go.uber.org/zap"
)

func With(k string, v interface{}) zap.Field {
	return zap.Any(k, v)
}

func WithError(err error) zap.Field {
	return zap.NamedError("error", err)
}

func Info(msg string, args ...zap.Field) {
	Log.Info(msg, args...)
}

func Error(msg string, args ...zap.Field) {
	Log.Error(msg, args...)
}

func Warn(msg string, args ...zap.Field) {
	Log.Warn(msg, args...)
}

func Debug(msg string, args ...zap.Field) {
	Log.Debug(msg, args...)
}

func Fatal(msg string, args ...zap.Field) {
	Log.Fatal(msg, args...)
}

func Infof(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	Log.Info(logMsg)
}

func Errorf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	Log.Error(logMsg)
}

func Warnf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	Log.Warn(logMsg)
}

func Debugf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	Log.Debug(logMsg)
}

func Fatalf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	Log.Fatal(logMsg)
}
