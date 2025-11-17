package tx

import "context"

type Tx struct {
	DbRepo DbRepo
}

func TxExecute(ctx context.Context, cb func(ctx context.Context) error) error {
	tx := fromContext(ctx)
	return tx.DbRepo.WithTx(ctx, cb)
}
