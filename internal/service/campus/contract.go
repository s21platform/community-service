package campus

import (
	"context"
	"github.com/s21platform/community-service/internal/model"
	"time"
)

type SchoolC interface {
	GetCampuses(ctx context.Context) ([]model.Campus, error)
}

type DbRepo interface {
	SetCampus(ctx context.Context, campus model.Campus) error
	GetCampusByUUID(ctx context.Context, campusUUID string) (*model.Campus, error)
}

type RedisRepo interface {
	GetByKey(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}
