package service

import "context"

type SchoolC interface {
	GetPeersByCampusUuid(ctx context.Context, campusUuid string, limit, offset int64) ([]string, error)
}

type DbRepo interface {
	AddPeerLogins(ctx context.Context, peerLogins []string) error
	GetCampusUuids(ctx context.Context) ([]string, error)
}
