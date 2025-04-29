package middleware

import (
	"time"

	"github.com/awend0/twitch-collector/internal/pkg/errcodes"
	"github.com/awend0/twitch-collector/internal/pkg/log"
	"github.com/ogen-go/ogen/middleware"
	"go.uber.org/zap"
)

func Logging(logger *zap.Logger) middleware.Middleware {
	return func(req middleware.Request, next func(req middleware.Request) (middleware.Response, error)) (middleware.Response, error) {
		lg := logger.With(
			zap.String("operationId", req.OperationID),
		)

		lg.Info("started request")

		req.Context = log.ToCtx(req.Context, lg)

		start := time.Now()
		resp, err := next(req)
		dur := time.Since(start).Microseconds()

		if err != nil {
			if errcode, ok := err.(*errcodes.ErrorCode); ok {
				lg.Error("request error",
					zap.Int("status_code", errcode.StatusCode),
					zap.String("message", errcode.Message),
					zap.String("details", errcode.Details),
				)
			} else {
				lg.Error("request internal error",
					zap.Error(err),
				)
			}
		}

		lg.Info("done request", zap.Float64("duration_ms", float64(dur)/float64(1000)))

		return resp, err
	}
}
