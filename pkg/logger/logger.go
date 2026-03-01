package logger

import (
	"context"

	"go.uber.org/zap"
)

const (
	key = "logger"
)

type Logger struct {
	l *zap.Logger
}

func New(ctx context.Context, env string) (context.Context, error) {
	lCfg := zap.NewDevelopmentConfig()
	switch env {
	case "dev":
		lCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "prod":
		lCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	l, err := lCfg.Build()
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, key, &Logger{l}), nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	if ctx.Value(key) == nil {
		return nil 
	}
	return ctx.Value(key).(*Logger)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.l.Fatal(msg, fields...)
}