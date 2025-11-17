package tx

import "context"

type DbRepo interface {
	WithTx(ctx context.Context, cb func(ctx context.Context) error) (err error)
}
