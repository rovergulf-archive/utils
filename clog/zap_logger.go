package clog

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

type (
	// LoggerOptions are options for constructing a Logger
	LoggerOptions struct {
		Level       string
		LogJSON     bool
		Development bool
		Stacktrace  bool
	}
	LogLevel string
)

func (l LogLevel) String() string {
	return string(l)
}

const (
	DebugLevel string = "DEBUG"
	InfoLevel  string = "INFO"
	WarnLevel  string = "WARN"
	ErrorLevel string = "ERROR"
	FatalLevel string = "FATAL"
)

var LogLevelMapping = map[string]zapcore.Level{
	DebugLevel: zap.DebugLevel,
	InfoLevel:  zap.InfoLevel,
	WarnLevel:  zap.WarnLevel,
	ErrorLevel: zap.ErrorLevel,
	FatalLevel: zap.FatalLevel,
}

// NewLogger creates a new Logger instance
func NewLogger(options LoggerOptions) (*zap.SugaredLogger, error) {
	config := zap.NewDevelopmentConfig()
	config.Development = options.Development
	config.DisableStacktrace = !options.Stacktrace

	if options.LogJSON {
		config.Encoding = "json"
	} else {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	options.Level = strings.ToUpper(options.Level)
	if _, exists := LogLevelMapping[options.Level]; !exists {
		options.Level = InfoLevel
	}

	config.Level = zap.NewAtomicLevelAt(LogLevelMapping[options.Level])

	logger, err := config.Build()
	if err != nil {
		logrus.Printf("Failed to enable logger: %s", err)
		return nil, err
	}

	l := logger.Sugar()

	l.Infow("Init Zap Logger.", "LEVEL", config.Level, "LOG_JSON", options.LogJSON)

	return l, nil
}

func init() {
	logrus.SetLevel(logrus.WarnLevel) // silence logs from zsais/go-gin-prometheus
}
