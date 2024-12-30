package service

import "context"

type SchoolS interface {
	GetPeersByCampusUuid(ctx context.Context, campusUuid string, limit, offset int64) ([]string, error)
}
