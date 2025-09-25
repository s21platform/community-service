package infra

import (
	"context"
	"net/http"

	logger_lib "github.com/s21platform/logger-lib"
	"google.golang.org/grpc"
)

func LoggerRPC(logger *logger_lib.Logger) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = logger_lib.NewContext(ctx, logger)
		return handler(ctx, req)
	}
}

func LoggerHTTP(logger *logger_lib.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = logger_lib.NewContext(ctx, logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
