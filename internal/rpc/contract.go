//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package rpc

import (
	"context"
	communityproto "github.com/s21platform/community-proto/community-proto"
	"github.com/s21platform/community-service/internal/model"
)

type DbRepo interface {
	SearchPeersBySubstring(ctx context.Context, substring string) ([]*communityproto.SearchPeer, error)
	GetPeerStatus(ctx context.Context, email string) (string, error)
	GetPeerSchoolData(ctx context.Context, nickName string) (model.PeerSchoolData, error)
}
