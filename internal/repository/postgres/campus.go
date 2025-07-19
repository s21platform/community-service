package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/s21platform/community-service/internal/model"
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

func (r *Repository) GetCampusByUUID(ctx context.Context, campusUUID string) (*model.Campus, error) {
	var campus model.Campus
	query, args, err := sq.Select("id", "campus_uuid").
		From("campus").
		Where(sq.Eq{"campus_uuid": campusUUID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to configure query, err: %v", err)
	}

	err = r.conn.GetContext(ctx, &campus, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get campus, err: %v", err)
	}

	return &campus, nil
}

func (r *Repository) SetCampus(ctx context.Context, campus model.Campus) error {
	query, args, err := sq.Insert("campus").Columns(
		`campus_uuid`,
		`short_name`,
		`full_name`,
	).Values(
		campus.Uuid,
		campus.ShortName,
		campus.FullName,
	).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to configure query, err: %v", err)
	}

	_, err = r.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert campus, err: %v", err)
	}

	return nil
}
