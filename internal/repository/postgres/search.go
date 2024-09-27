package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	communityproto "github.com/s21platform/community-proto/community-proto"
	"github.com/s21platform/community-service/internal/model"
	"log"
)

func (r *Repository) SearchPeersBySubstring(ctx context.Context, substring string) ([]*communityproto.SearchPeer, error) {
	var result []model.SearchPeers
	query, args, err := sq.Select(`login`).
		From("participant").
		Where("login LIKE ?", "%"+substring+"%").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	log.Println(query)
	//err := r.conn.Select(&result, `SELECT login FROM participant WHERE login LIKE $1`, "'%"+substring+"%'")
	if err != nil {
		return nil, fmt.Errorf("failed to get sql query: %w", err)
	}
	err = r.conn.Select(&result, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search peer by substring: %w", err)
	}
	log.Println(result)
	var res []*communityproto.SearchPeer
	for _, re := range result {
		res = append(res, &communityproto.SearchPeer{
			Login: re.Login,
		})
	}
	return res, nil
}
