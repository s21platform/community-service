package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) CheckLinkEduTwoPeers(ctx context.Context, uuidFirstPeer string, uuidSecondPeer string) (bool, error) {
	query, args, err := sq.Select("COUNT(DISTINCT id) = 2 AS has_both").
		From("link_edu").
		Where(sq.Eq{"id": []string{uuidFirstPeer, uuidSecondPeer}}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to CheckLinkEduTwoPeers, err: %v", err)
	}

	var flag bool
	err = r.conn.GetContext(ctx, &flag, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to CheckLinkEduTwoPeers, err: %v", err)
	}

	return flag, nil
}
