package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) GetIdPeer(ctx context.Context, uuid string) (int64, error) {
	var id int64
	query, args, err := sq.Select("edu_id").
		From("link_edu").
		Where(sq.Eq{"user_uuid": uuid}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build exists query: %v", err)
	}

	err = r.conn.GetContext(ctx, &id, query, args...)
	if err != nil {
		return 0, err
	}
	return id, nil
}
