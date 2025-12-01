package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Key string

const KeyTx = Key("tx_repo")

type Database interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func (r *Repository) Chk(ctx context.Context) Database {
	if v := getTx(ctx); v != nil {
		return v
	}
	return r.conn
}

func (r *Repository) WithTx(ctx context.Context, cb func(ctx context.Context) error) (err error) {
	if tx := getTx(ctx); tx != nil {
		return cb(ctx)
	}

	tx, err := r.conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	ctx = setTx(ctx, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				panic(err)
			}
		}
	}()

	err = cb(ctx)
	return err
}

func setTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, KeyTx, tx)
}

func getTx(ctx context.Context) *sqlx.Tx {
	v, ok := ctx.Value(KeyTx).(*sqlx.Tx)
	if !ok {
		return nil
	}
	return v
}
