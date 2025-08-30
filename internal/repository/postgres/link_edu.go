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

func (r *Repository) InsertLinkEdu(ctx context.Context, id int64, uuid string) error {
	query, args, err := sq.Insert("link_edu").
		Columns("edu_id", "user_uuid").
		Values(id, uuid).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %v", err)
	}

	_, err = r.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert link edu: %v", err)
	}

	return nil
}
