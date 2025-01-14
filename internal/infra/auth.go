package infra

import (
	"context"
	"database/sql"
	"errors"
	communityproto "github.com/s21platform/community-proto/community-proto"
	"github.com/s21platform/community-service/internal/repository/postgres"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/s21platform/community-service/internal/config"
)

func AuthInterceptor(env string, dbRepo postgres.Repository) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if info.FullMethod == "/CommunityService/IsPeerExist" {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)

		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "no info in metadata")
		}

		userIDs, ok := md["uuid"]
		if !ok || len(userIDs) != 1 {
			return nil, status.Errorf(codes.Unauthenticated, "no uuid or more than one in metadata")
		}
		ctx = context.WithValue(ctx, config.KeyUUID, userIDs[0])

		if env == "stage" {
			staffId, err := dbRepo.GetStaffIdByUserUuid(ctx, userIDs[0])
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return nil, status.Errorf(codes.Internal, "cannot check permission to the stage env, err: %v", err)
				}
			}

			if errors.Is(err, sql.ErrNoRows) || staffId <= 0 {
				return nil, status.Errorf(codes.PermissionDenied, "permission to the stage env denied")
			}
		}

		return handler(ctx, req)
	}
}
