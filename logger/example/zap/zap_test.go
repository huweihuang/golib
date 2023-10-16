package zap

import (
	"errors"
	"testing"
	"time"

	log "github.com/huweihuang/golib/logger/zap"
	"go.uber.org/zap"
)

func TestZap(t *testing.T) {
	c := log.New()
	c.SetDivision("size")  // 设置归档方式，"time"时间归档 "size" 文件大小归档，文件大小等可以在配置文件配置
	c.SetTimeUnit(log.Day) // 时间归档 可以设置切割单位
	c.SetEncoding("json")  // 输出格式 "json" 或者 "console"

	c.SetInfoFile("./log/zap.log")        // 设置日志文件
	c.SetErrorFile("./log/zap.error.log") // 设置error日志文件
	c.SetLogLevel("debug")
	c.InitLogger()

	PrintLog()
}

func TestZapLogByConfig(t *testing.T) {
	c := log.NewFromYaml("config.yaml")
	c.InitLogger()

	PrintLog()
}

func TestLog(t *testing.T) {
	log.Log().Infof("test default log")
	log.Log().Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", "example.com",
		"attempt", 3,
		"backoff", time.Second)

	log.Log().Error("test error log")
	log.Log().Errorw("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", "example.com",
		"attempt", 3,
		"backoff", time.Second)
}

func PrintLog() {
	// SugaredLogger
	log.SugaredLogger.Info("info level test")
	log.SugaredLogger.Error("error level test")
	log.SugaredLogger.Warn("warn level test")
	log.SugaredLogger.Debug("debug level test")

	log.SugaredLogger.Infof("info level test: %s", "111")
	log.SugaredLogger.Errorf("error level test: %s", "111")
	log.SugaredLogger.Warnf("warn level test: %s", "111")
	log.SugaredLogger.Debugf("debug level test: %s", "111")

	log.SugaredLogger.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", "example.com",
		"attempt", 3,
		"backoff", time.Second)

	// Logger
	log.Logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", "example.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
	log.Info("this is a log", log.With("Trace", "12345677"))
	log.Info("this is a log", log.WithError(errors.New("this is a new error")))
}