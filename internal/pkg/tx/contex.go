package tx

import (
	"context"

	"google.golang.org/grpc"
)

type Key string

const KeyTx = Key("tx")

func withDbRepoContext(ctx context.Context, repo DbRepo) context.Context {
	return context.WithValue(ctx, KeyTx, Tx{DbRepo: repo})
}

func TxMiddleWire(db DbRepo) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (
		interface{}, error) {
		return handler(withDbRepoContext(ctx, db), req)
	}
}

func fromContext(ctx context.Context) Tx {
	v, ok := ctx.Value(KeyTx).(Tx)
	if !ok {
		panic("no Tx found in context")
	}
	return v
}
