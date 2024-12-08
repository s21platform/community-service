package infra

import (
	"context"
	"github.com/s21platform/community-service/internal/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	_ = info
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no info in metadata")
	}

	userIDs, ok := md["uuid"]
	if !ok || len(userIDs) != 1 {
		return nil, status.Errorf(codes.Unauthenticated, "no uuid or more than one in metadata")
	}

	ctx = context.WithValue(ctx, config.KeyUUID, userIDs[0])

	return handler(ctx, req)
}
