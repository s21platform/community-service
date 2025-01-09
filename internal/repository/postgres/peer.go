package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) AddPeerLogins(ctx context.Context, peerLogins []string) error {
	queryBase := sq.Insert("login").
		Columns("nickname").
		Suffix("ON CONFLICT (nickname) DO NOTHING").
		PlaceholderFormat(sq.Dollar)

	for _, login := range peerLogins {
		queryBase = queryBase.Values(login)
	}

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
