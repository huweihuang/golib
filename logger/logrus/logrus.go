package logrus

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/huweihuang/golib/utils"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

const (
	defaultLevel             = "info"
	defaultLogFile           = "log/info.log"
	defaultFormat            = "json"
	defaultEnableForceColors = true
)

var (
	Logger *logrus.Logger

	callerPrettyfier = func(f *runtime.Frame) (string, string) {
		s := strings.Split(f.Function, ".")
		funcName := s[len(s)-1]
		fileName := path.Base(f.File)
		return funcName, fmt.Sprintf("%s:%d", fileName, f.Line)
	}
)

func InitLogger(logFile, logLevel, format string, enableForceColors bool) *logrus.Logger {
	logger := logrus.New()

	// set log level
	if logLevel == "" {
		logLevel = defaultLevel
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic("Failed to parse log level")
	}
	logger.SetLevel(level)

	// set stdout
	logger.SetOutput(os.Stdout)
	// set logfile if not empty
	if logFile != "" {
		accessLog := timeDivisionWriter(logFile)
		writers := []io.Writer{
			accessLog,
			os.Stdout,
		}
		fileAndStdoutWriter := io.MultiWriter(writers...)
		logger.SetOutput(fileAndStdoutWriter)
	}

	forceColors := false
	if enableForceColors {
		forceColors = true
	}
	// set file && line number
	logger.SetReportCaller(true)
	switch format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:  "2006-01-02 15:04:05",
			CallerPrettyfier: callerPrettyfier,
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:    true,
			ForceColors:      forceColors,
			DisableQuote:     true,
			TimestampFormat:  "2006-01-02 15:04:05",
			CallerPrettyfier: callerPrettyfier,
		})
	}

	Logger = logger
	return logger
}

func timeDivisionWriter(logFile string) io.Writer {
	err := utils.MakeParentDir(logFile)
	if err != nil {
		panic("Failed to create log directory")
	}

	accessLog, err := rotatelogs.New(
		logFile+".%Y%m%d",
		rotatelogs.WithLinkName(logFile),
		rotatelogs.WithMaxAge(time.Duration(7*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		panic("Failed to create access.log")
	}

	return accessLog
}

func InitDefaultLogger() *logrus.Logger {
	return InitLogger(defaultLogFile, defaultLevel, defaultFormat, defaultEnableForceColors)
}

// Log get a default Logger
func Log() *logrus.Logger {
	return InitDefaultLogger()
}
