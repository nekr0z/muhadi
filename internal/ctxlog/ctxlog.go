package ctxlog

import (
	"context"

	"go.uber.org/zap"
)

type key struct{}

func New(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, key{}, log)
}

func Maybe(ctx context.Context) *zap.Logger {
	log, ok := ctx.Value(key{}).(*zap.Logger)
	if !ok || log == nil {
		return zap.NewNop()
	}
	return log
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	log := Maybe(ctx)
	log.Error(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	log := Maybe(ctx)
	log.Warn(msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	log := Maybe(ctx)
	log.Info(msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	log := Maybe(ctx)
	log.Debug(msg, fields...)
}
