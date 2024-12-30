package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) GetCampusUuids(ctx context.Context) ([]string, error) {
	var campuses []string
	query, args, err := sq.Select("campus_uuid").
		From("campus").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot configure query, err: %v", err)
	}

	err = r.conn.SelectContext(ctx, &campuses, query, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot get campus, err: %v", err)
	}

	return campuses, nil
}
