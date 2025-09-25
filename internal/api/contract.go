//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package api

import (
	"context"
	"time"

	"github.com/s21platform/community-service/internal/config"
)

type DbRepo interface {
	GetPeerStatus(ctx context.Context, login string) (string, error)
	GetIdFromParticipant(ctx context.Context, login string) (int64, error)
	InsertLinkEdu(ctx context.Context, id int64, uuid string) error
}

type RedisRepo interface {
	Set(ctx context.Context, key config.Key, value string, expiration time.Duration) error
	GetByKey(ctx context.Context, key config.Key) (string, error)
}

type NotificationClient interface {
	SendEduCode(ctx context.Context, email, code string) error
}
