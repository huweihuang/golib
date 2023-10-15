package zap

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/huweihuang/golib/utils"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	TimeDivision = "time"
	SizeDivision = "size"

	defaultEncoding      = "console"
	defaultDivision      = "size"
	defaultUnit          = Day
	defaultLogLevel      = "info"
	defaultInfoFilename  = "log/info.log"
	defaultErrorFilename = "log/error.log"

	// lumberjack
	defaultMaxSize    = 256 //MB
	defaultMaxBackups = 30
	defaultMaxAge     = 30
)

var (
	SugaredLogger *zap.SugaredLogger
	Logger        *zap.Logger

	_encoderNameToConstructor = map[string]func(zapcore.EncoderConfig) zapcore.Encoder{
		"console": func(encoderConfig zapcore.EncoderConfig) zapcore.Encoder {
			return zapcore.NewConsoleEncoder(encoderConfig)
		},
		"json": func(encoderConfig zapcore.EncoderConfig) zapcore.Encoder {
			return zapcore.NewJSONEncoder(encoderConfig)
		},
	}
)

type LogOptions struct {
	// Encoding sets the logger's encoding. Valid values are "json" and
	// "console", as well as any third-party encodings registered via
	// RegisterEncoder.
	Encoding      string   `json:"encoding" yaml:"encoding" toml:"encoding"`
	InfoFilename  string   `json:"info_filename" yaml:"info_filename" toml:"info_filename"`
	ErrorFilename string   `json:"error_filename" yaml:"error_filename" toml:"error_filename"`
	MaxSize       int      `json:"max_size" yaml:"max_size" toml:"max_size"`
	MaxBackups    int      `json:"max_backups" yaml:"max_backups" toml:"max_backups"`
	MaxAge        int      `json:"max_age" yaml:"max_age" toml:"max_age"`
	Compress      bool     `json:"compress" yaml:"compress" toml:"compress"`
	Division      string   `json:"division" yaml:"division" toml:"division"`
	LevelSeparate bool     `json:"level_separate" yaml:"level_separate" toml:"level_separate"`
	TimeUnit      TimeUnit `json:"time_unit" yaml:"time_unit" toml:"time_unit"`
	LogLevel      string   `json:"log_level" yaml:"log_level" toml:"log_level"`

	consoleDisplay bool
	caller         bool
}

func New() *LogOptions {
	return &LogOptions{
		Encoding:       defaultEncoding,
		InfoFilename:   defaultInfoFilename,
		ErrorFilename:  defaultErrorFilename,
		Division:       defaultDivision,
		LevelSeparate:  false,
		LogLevel:       defaultLogLevel,
		TimeUnit:       defaultUnit,
		MaxSize:        defaultMaxSize, //MB
		MaxBackups:     defaultMaxBackups,
		MaxAge:         defaultMaxAge, //days
		Compress:       true,
		caller:         true,
		consoleDisplay: true,
	}
}

// Log get a default SugaredLogger
func Log() *zap.SugaredLogger {
	log := New()
	log.InitLogger()
	return SugaredLogger
}

// NewLog new a SugaredLogger by logFile, logLevel, format
func NewLog(logFile, logLevel, format string) *zap.SugaredLogger {
	log := New()
	log.SetInfoFile(logFile)
	log.SetLogLevel(logLevel)
	log.SetEncoding(format)
	log.InitLogger()
	return SugaredLogger
}

func (c *LogOptions) InitLogger() (*zap.Logger, *zap.SugaredLogger) {
	var (
		core               zapcore.Core
		infoHook, warnHook io.Writer
		wsInfo             []zapcore.WriteSyncer
		wsWarn             []zapcore.WriteSyncer
	)

	if c.Encoding == "" {
		c.Encoding = defaultEncoding
	}
	encoder := _encoderNameToConstructor[c.Encoding]

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if c.consoleDisplay {
		wsInfo = append(wsInfo, zapcore.AddSync(os.Stdout))
		wsWarn = append(wsWarn, zapcore.AddSync(os.Stdout))
	}

	// zapcore WriteSyncer setting
	if c.isOutput() {
		switch c.Division {
		case TimeDivision:
			infoHook = c.timeDivisionWriter(c.InfoFilename)
			if c.LevelSeparate {
				warnHook = c.timeDivisionWriter(c.ErrorFilename)
			}
		case SizeDivision:
			infoHook = c.sizeDivisionWriter(c.InfoFilename)
			if c.LevelSeparate {
				warnHook = c.sizeDivisionWriter(c.ErrorFilename)
			}
		}
		wsInfo = append(wsInfo, zapcore.AddSync(infoHook))
	}

	if c.ErrorFilename != "" {
		wsWarn = append(wsWarn, zapcore.AddSync(warnHook))
	}

	// Separate info and warning log
	if c.LevelSeparate {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder(encoderConfig), zapcore.NewMultiWriteSyncer(wsInfo...), infoLevel()),
			zapcore.NewCore(encoder(encoderConfig), zapcore.NewMultiWriteSyncer(wsWarn...), warnLevel()),
		)
	} else {
		level := convertLogLevel(c.LogLevel)
		core = zapcore.NewCore(encoder(encoderConfig), zapcore.NewMultiWriteSyncer(wsInfo...), level)
	}

	// file line number display
	development := zap.Development()
	// init default key
	//filed := zap.Fields(zap.String("serviceName", "serviceName"))
	var logger *zap.Logger
	if c.caller {
		logger = zap.New(core, zap.AddCaller(), development)
	} else {
		logger = zap.New(core, development)
	}

	Logger = logger
	SugaredLogger = Logger.Sugar()
	return Logger, SugaredLogger
}

func (c *LogOptions) SetDivision(division string) {
	c.Division = division
}

func (c *LogOptions) CloseConsoleDisplay() {
	c.consoleDisplay = false
}

func (c *LogOptions) defaultDisplay() {
	c.consoleDisplay = true
}

func (c *LogOptions) SetCaller(b bool) {
	c.caller = b
}

func (c *LogOptions) SetTimeUnit(t TimeUnit) {
	c.TimeUnit = t
}

func (c *LogOptions) SetErrorFile(path string) {
	c.LevelSeparate = true
	c.ErrorFilename = path
}

func (c *LogOptions) SetInfoFile(path string) {
	c.InfoFilename = path
}

func (c *LogOptions) SetEncoding(encoding string) {
	c.Encoding = encoding
}

func (c *LogOptions) SetLogLevel(level string) {
	c.LogLevel = level
}

// isOutput whether set output file
func (c *LogOptions) isOutput() bool {
	return c.InfoFilename != ""
}

func (c *LogOptions) sizeDivisionWriter(filename string) io.Writer {
	err := utils.MakeParentDir(filename)
	if err != nil {
		panic("Failed to create log directory")
	}

	hook := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
		Compress:   c.Compress,
		LocalTime:  true,
	}
	return hook
}

func (c *LogOptions) timeDivisionWriter(filename string) io.Writer {
	err := utils.MakeParentDir(filename)
	if err != nil {
		panic("Failed to create log directory")
	}

	hook, err := rotatelogs.New(
		filename+c.TimeUnit.Format(),
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Duration(int64(24*time.Hour)*int64(c.MaxAge))),
		rotatelogs.WithRotationTime(c.TimeUnit.RotationGap()),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func Info(msg string, args ...zap.Field) {
	Logger.Info(msg, args...)
}

func Error(msg string, args ...zap.Field) {
	Logger.Error(msg, args...)
}

func Warn(msg string, args ...zap.Field) {
	Logger.Warn(msg, args...)
}

func Debug(msg string, args ...zap.Field) {
	Logger.Debug(msg, args...)
}

func Fatal(msg string, args ...zap.Field) {
	Logger.Fatal(msg, args...)
}

func Infof(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	Logger.Info(logMsg)
}

func Errorf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	Logger.Error(logMsg)
}

func Warnf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	Logger.Warn(logMsg)
}

func Debugf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	Logger.Debug(logMsg)
}

func Fatalf(format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)
	Logger.Fatal(logMsg)
}

func With(k string, v interface{}) zap.Field {
	return zap.Any(k, v)
}

func WithError(err error) zap.Field {
	return zap.NamedError("error", err)
}

func infoLevel() zap.LevelEnablerFunc {
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})
}

func warnLevel() zap.LevelEnablerFunc {
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})
}

func convertLogLevel(logLevel string) (level zapcore.Level) {
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	return level
}
