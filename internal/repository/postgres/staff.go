package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) GetStaffId(ctx context.Context, login string) (int64, error) {
	query, args, err := sq.Select("id").
		From("staff").
		Where(sq.Eq{"login": login}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to configure query, err: %v", err)
	}

	var id int64
	err = r.conn.GetContext(ctx, &id, query, args...)
	if err != nil {
		return 0, err
	}

	return id, nil
}
