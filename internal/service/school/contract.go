package service

import (
	"context"
	"time"

	"github.com/s21platform/community-service/internal/config"
)

type SchoolClient interface {
	GetPeersByCampusUuid(ctx context.Context, campusUuid string, limit, offset int64) ([]string, error)
}

type DbRepo interface {
	AddPeerLogins(ctx context.Context, peerLogins []string) error
	GetCampusUuids(ctx context.Context) ([]string, error)
}

type RedisRepo interface {
	GetByKey(ctx context.Context, key config.Key) (string, error)
	Set(ctx context.Context, key config.Key, value string, expiration time.Duration) error
}
