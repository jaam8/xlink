package logger

import (
	"context"
	"go.uber.org/zap"
)

type key string

const (
	KeyForLogger    key = "logger"
	KeyForRequestID key = "request_id"
)

type Logger struct {
	l *zap.Logger
}

func NewLogger() (*Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	loggerStruct := &Logger{l: logger}

	return loggerStruct, nil
}

func New(ctx context.Context) (context.Context, error) {
	loggerStruct, err := NewLogger()
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, KeyForLogger, loggerStruct)
	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	return ctx.Value(KeyForLogger).(*Logger)
}

func TryAppendRequestIDFromContext(ctx context.Context, fields []zap.Field) []zap.Field {
	if ctx.Value(KeyForRequestID) != nil {
		fields = append(fields, zap.String(string(KeyForRequestID), ctx.Value(KeyForRequestID).(string)))
	}
	return fields
}

func GetOrCreateLoggerFromCtx(ctx context.Context) *Logger {
	logger := GetLoggerFromCtx(ctx)
	if logger == nil {
		logger, _ = NewLogger()
	}
	return logger
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	fields = TryAppendRequestIDFromContext(ctx, fields)
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = TryAppendRequestIDFromContext(ctx, fields)
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	fields = TryAppendRequestIDFromContext(ctx, fields)
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = TryAppendRequestIDFromContext(ctx, fields)
	l.l.Error(msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	fields = TryAppendRequestIDFromContext(ctx, fields)
	l.l.Fatal(msg, fields...)
}
