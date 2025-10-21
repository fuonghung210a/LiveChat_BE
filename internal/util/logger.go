package util

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevel represents the logging level
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
)

// NewLogger creates a new logger with the specified log level
// If no level is provided or invalid, defaults to INFO level
func NewLogger(level ...LogLevel) *zap.Logger {
	var logLevel LogLevel
	if len(level) > 0 {
		logLevel = level[0]
	} else {
		// Get log level from environment variable
		envLevel := os.Getenv("LOG_LEVEL")
		if envLevel != "" {
			logLevel = LogLevel(envLevel)
		} else {
			logLevel = InfoLevel
		}
	}

	// Convert LogLevel to zapcore.Level
	var zapLevel zapcore.Level
	switch logLevel {
	case DebugLevel:
		zapLevel = zapcore.DebugLevel
	case InfoLevel:
		zapLevel = zapcore.InfoLevel
	case WarnLevel:
		zapLevel = zapcore.WarnLevel
	case ErrorLevel:
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Configure encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Create core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	// Build logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return logger
}

// GetLogger returns a logger instance (can be used to get a shared logger)
func GetLogger() *zap.Logger {
	return NewLogger()
}

