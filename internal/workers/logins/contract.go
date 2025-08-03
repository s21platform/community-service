package logins

import (
	"context"
	"time"

	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/model"
)

type SchoolClient interface {
	GetPeersByCampusUuid(ctx context.Context, campusUuid string, limit, offset int64) ([]string, error)
}

type DbRepo interface {
	GetPeerByLogin(ctx context.Context, nickname string) (model.Login, error)
	SetNickname(ctx context.Context, nickname string) error
	GetCampusUuids(ctx context.Context) ([]string, error)
}

type RedisRepo interface {
	GetByKey(ctx context.Context, key config.Key) (string, error)
	Set(ctx context.Context, key config.Key, value string, expiration time.Duration) error
}
