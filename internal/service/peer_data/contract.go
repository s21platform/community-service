package service

import (
	"context"
	"time"

	school "github.com/s21platform/school-proto/school-proto"
)

type SchoolC interface {
	GetPeersByCampusUuid(ctx context.Context, campusUuid string, limit, offset int64) ([]string, error)
	GetParticipantData(ctx context.Context, login string) (*school.GetParticipantDataOut, error)
}

type DbRepo interface {
	AddPeerLogins(ctx context.Context, peerLogins []string) error
	GetCampusUuids(ctx context.Context) ([]string, error)
	SaveParticipantData(ctx context.Context, data *school.GetParticipantDataOut) error
}

type RedisRepo interface {
	GetByKey(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}