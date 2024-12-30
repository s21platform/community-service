package service

import (
	"context"
)

type DbRepo interface {
	AddPeerLogins(ctx context.Context, peerLogins []string) error
	//лучше uuid кампусов передавать как слайс строки или создать структуру кампуса с полем uuid и передавать слайс этих структур?
	GetCampusUuids(ctx context.Context) ([]string, error)
}
