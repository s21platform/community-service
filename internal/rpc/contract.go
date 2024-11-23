package rpc

import (
	"context"
	communityproto "github.com/s21platform/community-proto/community-proto"
)

type DbRepo interface {
	SearchPeersBySubstring(ctx context.Context, substring string) ([]*communityproto.SearchPeer, error)
	GetPeerStatus(ctx context.Context, email string) (string, error)
}
