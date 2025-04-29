package log

import (
	"context"

	"go.uber.org/zap"
)

type ctxLogKey struct{}

func ToCtx(ctx context.Context, lg *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLogKey{}, lg)
}

func Ctx(ctx context.Context) *zap.Logger {
	lg, ok := ctx.Value(ctxLogKey{}).(*zap.Logger)
	if !ok || lg == nil {
		zap.L().Error("no logger in context, using global")
		return zap.L()
	}

	return lg
}

func Global() *zap.Logger {
	return zap.L()
}
