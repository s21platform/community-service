//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"

	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/model"
	"github.com/s21platform/community-service/pkg/community"
)

type DbRepo interface {
	Conn() *sqlx.DB
	SearchPeersBySubstring(ctx context.Context, substring string) ([]*community.SearchPeer, error)
	GetPeerStatus(ctx context.Context, login string) (string, error)
	GetPeerSchoolData(ctx context.Context, nickName string) (model.PeerSchoolData, error)
	GetStaffId(ctx context.Context, login string) (int64, error)
	GetPeerData(ctx context.Context, id int64) (*model.ParticipantData, error)
	GetIdPeer(ctx context.Context, uuid string) (int64, error)
	GetIdFromParticipant(ctx context.Context, login string) (int64, error)
	InsertLinkEdu(ctx context.Context, id int64, uuid string) error
	GetLogin(ctx context.Context, id int64) (string, error)
}

type RedisRepo interface {
	GetByKey(ctx context.Context, key config.Key) (string, error)
	Set(ctx context.Context, key config.Key, value string, expiration time.Duration) error
	Delete(ctx context.Context, key config.Key)
}

type NotificationS interface {
	SendEduCode(ctx context.Context, email, code string) error
}

type UserLinkingEdu interface {
	ProduceMessage(ctx context.Context, message any, key any) error
}
