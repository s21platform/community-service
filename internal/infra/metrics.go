package infra

import (
	"context"
	"github.com/s21platform/community-service/internal/config"
	"strings"
	"time"

	"github.com/s21platform/metrics-lib/pkg"
	"google.golang.org/grpc"
)

func MetricsInterceptor(metrics *pkg.Metrics) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (
		interface{}, error) {
		t := time.Now()
		method := strings.Trim(strings.ReplaceAll(info.FullMethod, "/", "_"), "_")
		metrics.Increment(method)

		ctx = context.WithValue(ctx, config.KeyMetrics, metrics)
		resp, err := handler(ctx, req)

		if err != nil {
			metrics.Increment(method + "_error")
		}

		metrics.Duration(time.Since(t).Milliseconds(), method)

		return resp, err
	}
}
