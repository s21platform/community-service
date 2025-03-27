package postgres

import (
	"context"
	// "database/sql"
	// "errors"
	"fmt"
	// `github.com/s21platform/community-service/internal/model`
	school "github.com/s21platform/school-proto/school-proto"
	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) GetParticinatsLogin(ctx context.Context) ([]string, error) {
	var loginsParticipants []string
	query, args, err := sq.Select("campus_uuid").
		From("campus").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot configure query, err: %v", err)
	}

	err = r.conn.SelectContext(ctx, &loginsParticipants, query, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot get campus, err: %v", err)
	}

	return loginsParticipants, nil
}


func (r *Repository) SetParticipantData(ctx context.Context, participantData *school.GetParticipantDataOut ) error {
	queryBase := sq.Insert("login").
		Columns("nickname").
		Suffix("ON CONFLICT (nickname) DO NOTHING").
		PlaceholderFormat(sq.Dollar)

	// for _, login := range participantData  {
	// 	queryBase = queryBase.Values(login)
	// }

	query, args, err := queryBase.ToSql()
	if err != nil {
		return fmt.Errorf("cannot configure query, err: %v", err)
	}
	_, err = r.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("cannot execute query, err: %v", err)
	}

	return nil
}
