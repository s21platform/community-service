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
		return 0, fmt.Errorf("cannot configure query, err: %v", err)
	}

	var id int64
	err = r.conn.GetContext(ctx, &id, query, args...)
	if err != nil {
		return 0, fmt.Errorf("cannot get staff id, err: %v", err)
	}

	return id, nil
}

// TODO: понять как проверить является ли юзер стафом
func (r *Repository) GetStaffIdByUserUuid(ctx context.Context, userUuid string) (int64, error) {
	//query, args, err := sq.Select("id").
	//	From("staff").
	//	Join("participant on staff.login = participant.login").
	//	Where(sq.Eq{"participant.uuid"})
	return 0, nil
}
