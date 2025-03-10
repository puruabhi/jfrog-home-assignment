package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	*zap.SugaredLogger
}

func NewZapLogger() *zapLogger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel) // Set default log level to info

	// Change the encoder configuration to human-readable format
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Encoding = "console"

	logger, _ := config.Build()
	defer logger.Sync()

	sugar := logger.Sugar()
	return &zapLogger{SugaredLogger: sugar}
}
