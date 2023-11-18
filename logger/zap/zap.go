package zap

import (
	"io"
	"os"
	"time"

	"github.com/huweihuang/golib/utils"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	TimeDivision = "time"
	SizeDivision = "size"
	TextFormat   = "text"
	JsonFormat   = "json"
)

const (
	defaultInfoFilename  = "log/info.log"
	defaultErrorFilename = "log/error.log"
	defaultLogLevel      = "info"
	defaultEncoding      = "json"

	defaultDivision = TimeDivision
	defaultUnit     = Day

	// lumberjack
	defaultMaxSize    = 256 //MB
	defaultMaxBackups = 30
	defaultMaxAge     = 30
)

var (
	logFilled atomic.Bool

	SugaredLogger *zap.SugaredLogger
	L             *zap.Logger

	_encoderNameToConstructor = map[string]func(zapcore.EncoderConfig) zapcore.Encoder{
		TextFormat: func(encoderConfig zapcore.EncoderConfig) zapcore.Encoder {
			return zapcore.NewConsoleEncoder(encoderConfig)
		},
		JsonFormat: func(encoderConfig zapcore.EncoderConfig) zapcore.Encoder {
			return zapcore.NewJSONEncoder(encoderConfig)
		},
	}
)

type LogOptions struct {
	// Encoding sets the logger's encoding. Valid values are "json" and
	// "text", as well as any third-party encodings registered via
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
	EnableColor   bool     `json:"enable_color" yaml:"enable_color" toml:"enable_color"`

	consoleOutput bool
	caller        bool
}

func New() *LogOptions {
	return &LogOptions{
		Encoding:      defaultEncoding,
		InfoFilename:  "",
		ErrorFilename: "",
		Division:      defaultDivision,
		LevelSeparate: false,
		LogLevel:      defaultLogLevel,
		TimeUnit:      defaultUnit,
		MaxSize:       defaultMaxSize, //MB
		MaxBackups:    defaultMaxBackups,
		MaxAge:        defaultMaxAge, //days
		Compress:      true,
		caller:        true,
		consoleOutput: true,
		EnableColor:   false,
	}
}

// Log get a default L
func Log() *zap.Logger {
	if logFilled.Load() {
		// return an initialized logger
		return L
	}
	log := New()
	log.InitLogger()
	return L
}

// Logger get a default SlgaredLogger
func Logger() *zap.SugaredLogger {
	if logFilled.Load() {
		// return an initialized logger
		return SugaredLogger
	}
	log := New()
	log.InitLogger()
	return SugaredLogger
}

// Sugar is the alias of Logger
func Sugar() *zap.SugaredLogger {
	return Logger()
}

// InitLogger new a Logger and SugaredLogger by logFile, logLevel, format
func InitLogger(logFile, errFile, logLevel, format string, enableColor bool) (*zap.Logger, *zap.SugaredLogger) {
	log := New()
	log.SetInfoFile(logFile)
	log.SetErrorFile(errFile)
	log.SetLogLevel(logLevel)
	log.SetEncoding(format)
	log.SetColor(enableColor)
	log.InitLogger()
	return L, SugaredLogger
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

	if c.EnableColor {
		encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	}

	if c.consoleOutput {
		wsInfo = append(wsInfo, zapcore.AddSync(os.Stdout))
		wsWarn = append(wsWarn, zapcore.AddSync(os.Stdout))
	}

	// zapcore WriteSyncer setting
	if c.InfoFilename != "" {
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

	L = logger
	SugaredLogger = L.Sugar()

	// mark as already initialized
	logFilled.Store(true)
	return L, SugaredLogger
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

func (c *LogOptions) SetDivision(division string) {
	c.Division = division
}

func (c *LogOptions) SetTimeUnit(t TimeUnit) {
	c.TimeUnit = t
}

func (c *LogOptions) SetConsoleDisplay(flag bool) {
	c.consoleOutput = flag
}

func (c *LogOptions) SetCaller(flag bool) {
	c.caller = flag
}

func (c *LogOptions) SetColor(flag bool) {
	c.EnableColor = flag
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
