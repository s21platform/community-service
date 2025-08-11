//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"context"

	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/model"
	"github.com/s21platform/community-service/pkg/community"
)

type DbRepo interface {
	SearchPeersBySubstring(ctx context.Context, substring string) ([]*community.SearchPeer, error)
	GetPeerStatus(ctx context.Context, email string) (string, error)
	GetPeerSchoolData(ctx context.Context, nickName string) (model.PeerSchoolData, error)
	GetStaffId(ctx context.Context, login string) (int64, error)
}

type RedisRepo interface {
	Delete(ctx context.Context, key config.Key)
}

type LevelChangeProducer interface {
	ProduceMessage(ctx context.Context, message any, key any) error
}

type producerExpLevelChanged interface {
	ProduceMessage(ctx context.Context, message any, key any) error
}

type producerStatusChanged interface {
	ProduceMessage(ctx context.Context, message any, key any) error
}