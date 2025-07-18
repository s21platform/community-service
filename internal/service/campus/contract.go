package campus

import (
	"context"
	"time"

	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/model"
)

type SchoolClient interface {
	GetCampuses(ctx context.Context) ([]model.Campus, error)
}

type DbRepo interface {
	SetCampus(ctx context.Context, campus model.Campus) error
	GetCampusByUUID(ctx context.Context, campusUUID string) (*model.Campus, error)
}

type RedisRepo interface {
	GetByKey(ctx context.Context, key config.Key) (string, error)
	Set(ctx context.Context, key config.Key, value string, expiration time.Duration) error
}
