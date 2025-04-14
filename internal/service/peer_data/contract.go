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
	GetParticipantData(ctx context.Context, participantData *school.GetParticipantDataOut )(*school.GetParticipantDataOut, error)
	GetParticipantsLogin(ctx context.Context, limit, offset int64) ([]string, error)
	SaveParticipantData(ctx context.Context, participantData *school.GetParticipantDataOut, login string) error
}

type RedisRepo interface {
	GetByKey(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}