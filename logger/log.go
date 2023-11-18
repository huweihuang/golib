package logger

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	logrus2 "github.com/huweihuang/golib/logger/logrus"
	zap2 "github.com/huweihuang/golib/logger/zap"
)

const (
	ZapType    = "zap"
	LogrusType = "logrus"
)

var (
	zapLogFilled    atomic.Bool
	logrusLogFilled atomic.Bool

	ZapSugaredLogger *zap.SugaredLogger
	ZapLogger        *zap.Logger
	LogrusLogger     *logrus.Logger
)

type LogConfig struct {
	LogFile     string
	ErrorFile   string
	LogLevel    string
	LogFormat   string
	EnableColor bool
}

func InitLogger(c *LogConfig, loggerType string) (*zap.Logger, *zap.SugaredLogger, *logrus.Logger) {
	switch loggerType {
	case ZapType:
		ZapLogger, ZapSugaredLogger = zap2.InitLogger(c.LogFile, c.ErrorFile, c.LogLevel, c.LogFormat, c.EnableColor)
		zapLogFilled.Store(true)
	case LogrusType:
		LogrusLogger = logrus2.InitLogger(c.LogFile, c.LogLevel, c.LogFormat, c.EnableColor)
		logrusLogFilled.Store(true)
	}
	return ZapLogger, ZapSugaredLogger, LogrusLogger
}

func Logrus() *logrus.Logger {
	if logrusLogFilled.Load() {
		return LogrusLogger
	}
	return logrus2.InitDefaultLogger()
}

func Sugar() *zap.SugaredLogger {
	if zapLogFilled.Load() {
		return ZapSugaredLogger
	}
	return zap2.Logger()
}

func Zap() *zap.Logger {
	if zapLogFilled.Load() {
		return ZapLogger
	}
	return zap2.Log()
}
