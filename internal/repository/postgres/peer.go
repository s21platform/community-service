package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/s21platform/community-service/internal/model"
)

func (r *Repository) GetPeerByLogin(ctx context.Context, nickname string) (model.Login, error) {
	query, args, err := sq.Select("login").
		Column("nickname").
		Where(sq.Eq{"nickname": nickname}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return model.Login{}, fmt.Errorf("failed to build query: %v", err)
	}
	var result model.Login
	err = r.conn.GetContext(ctx, &result, query, args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Login{}, nil
		}
		return model.Login{}, fmt.Errorf("failed to run query: %v", err)
	}
	return result, nil
}

func (r *Repository) SetNickname(ctx context.Context, nickname string) error {
	query, args, err := sq.Insert("login").
		Columns("nickname").
		Values(nickname).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}
	_, err = r.conn.ExecContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed to run query: %v", err)
	}

	return nil
}
