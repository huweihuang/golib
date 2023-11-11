package logger

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"

	logrus2 "github.com/huweihuang/golib/logger/logrus"
	zap2 "github.com/huweihuang/golib/logger/zap"
)

const (
	ZapType    = "zap"
	LogrusType = "logrus"
)

var (
	ZapLogger    *zap.SugaredLogger
	LogrusLogger *logrus.Logger
)

type LogConfig struct {
	LogFile           string
	LogLevel          string
	LogFormat         string
	EnableForceColors bool
}

func InitLogger(c *LogConfig, loggerType string) (*zap.SugaredLogger, *logrus.Logger) {
	switch loggerType {
	case ZapType:
		ZapLogger = zap2.InitLogger(c.LogFile, c.LogLevel, c.LogFormat)
	case LogrusType:
		LogrusLogger = logrus2.InitLogger(c.LogFile, c.LogLevel, c.LogFormat, c.EnableForceColors)
	}
	return ZapLogger, LogrusLogger
}
