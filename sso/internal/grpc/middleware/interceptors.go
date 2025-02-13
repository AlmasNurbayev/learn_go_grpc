package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ctxKey string

const TraceIDKey ctxKey = "traceID"

func ContextMiddleware(timeout time.Duration, log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		const op = "middleware.ContextMiddleware"
		log.With(slog.String("op", op))

		ctx = context.WithValue(ctx, TraceIDKey, "yes")
		ctx, cancel := context.WithTimeout(ctx, timeout)

		defer cancel()

		start := time.Now()
		fmt.Println("middleware", ctx.Value(TraceIDKey))
		res, err := handler(ctx, req)
		elapsed := time.Since(start)
		log.Info("request duration", slog.String("method", info.FullMethod), slog.Duration("duration", elapsed))

		if ctx.Err() == context.DeadlineExceeded {
			log.Warn("Timeout exceeded")
			return nil, status.Error(codes.DeadlineExceeded, "Timeout exceeded")
		}
		if ctx.Err() == context.Canceled {
			log.Warn("Context canceled")
			return nil, status.Error(codes.Canceled, "Context canceled")
		}

		return res, err
	}
}
