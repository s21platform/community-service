package service

import "context"

type SchoolC interface {
	GetPeersByCampusUuid(ctx context.Context, campusUuid string, limit, offset int64) ([]string, error)
}
