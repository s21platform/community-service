//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package api

import (
	"context"
	"github.com/s21platform/community-service/internal/config"
	"time"
)

type DbRepo interface {
	GetPeerStatus(ctx context.Context, login string) (string, error)
}

type RedisRepo interface {
	Set(ctx context.Context, key config.Key, value string, expiration time.Duration) error
}

type NotificationClient interface {
	SendEduCode(ctx context.Context, email, code string) error
}
