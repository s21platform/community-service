package postgres

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/s21platform/community-service/pkg/community"

	"github.com/s21platform/community-service/internal/model"
)

func (r *Repository) SearchPeersBySubstring(ctx context.Context, substring string) ([]*community.SearchPeer, error) {
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
	var res []*community.SearchPeer
	for _, re := range result {
		res = append(res, &community.SearchPeer{
			Login: re.Login,
		})
	}
	return res, nil
}
