package service

import (
	"context"
	"time"
)

type SchoolC interface {
	GetPeersByCampusUuid(ctx context.Context, campusUuid string, limit, offset int64) ([]string, error)
}

type DbRepo interface {
	AddPeerLogins(ctx context.Context, peerLogins []string) error
	GetCampusUuids(ctx context.Context) ([]string, error)
}

type RedisRepo interface {
	Get(ctx context.Context) (string, error)
	GetByKey(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}
