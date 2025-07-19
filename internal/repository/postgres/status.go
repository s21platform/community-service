package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) GetPeerStatus(ctx context.Context, email string) (string, error) {
	var status string
	query, args, err := sq.Select("status").
		From("participant").
		Where(sq.Eq{"login": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("failed to configure query, err: %v", err)
	}

	err = r.conn.GetContext(ctx, &status, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return status, nil
}
