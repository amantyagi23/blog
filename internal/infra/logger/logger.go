package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap for structured logging.
type Logger struct {
	*zap.Logger
}

// New creates a production-ready logger.
func New(env string) (*Logger, error) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return &Logger{logger}, nil
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() {
	_ = l.Logger.Sync()
}

// WithContext adds context fields to logger.
func (l *Logger) WithContext(fields ...zap.Field) *Logger {
	return &Logger{l.Logger.With(fields...)}
}