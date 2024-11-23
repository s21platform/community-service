package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) IsPeerExist(ctx context.Context, email string) (string, error) {
	var status string
	query, args, err := sq.Select("status").
		From("participant").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("cannot configure query, err: %v", err)
	}

	err = r.conn.QueryRowContext(ctx, query, args...).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("cannot get peer status, err: %v", err)
	}
	return status, nil
}
