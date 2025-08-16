package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) CheckLinkEduTwoPeers(ctx context.Context, uuidFirstPeer string, uuidSecondPeer string) (int64, error) {
	query, args, err := sq.Select("COUNT(DISTINCT user_uuid)").
		From("link_edu").
		Where(sq.Eq{"user_uuid": []string{uuidFirstPeer, uuidSecondPeer}}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to Check Link Edu Two Peers, err: %v", err)
	}

	var flag int64
	err = r.conn.GetContext(ctx, &flag, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to Check Link Edu Two Peers, err: %v", err)
	}

	return flag, nil
}

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
