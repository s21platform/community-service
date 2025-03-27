package service

import (
	"context"
	"time"

	school "github.com/s21platform/school-proto/school-proto"
)

type SchoolC interface {
	GetParticipantData(ctx context.Context, login string) (*school.GetParticipantDataOut, error)
}

type DbRepo interface {
	SetParticipantData(ctx context.Context, participantData *school.GetParticipantDataOut )(error)
	GetParticinatsLogin(ctx context.Context) ([]string, error)
	SaveParticipantData(ctx context.Context, data *school.GetParticipantDataOut) error
}

type RedisRepo interface {
	GetByKey(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}