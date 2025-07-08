package peerdata

import (
	"context"
	"time"

	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/model"
)

type SchoolC interface {
	GetParticipantData(ctx context.Context, login string) (*model.ParticipantDataValue, error)
}

type DbRepo interface {
	GetParticipantsLogin(ctx context.Context, limit, offset int64) ([]string, error)
	SetParticipantData(ctx context.Context, participantDataValue *model.ParticipantDataValue, login string) error
}

type RedisRepo interface {
	GetByKey(ctx context.Context, key config.Key) (string, error)
	Set(ctx context.Context, key config.Key, value string, expiration time.Duration) error
}
