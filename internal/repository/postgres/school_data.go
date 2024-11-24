package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/s21platform/community-service/internal/model"
)

func (r *Repository) GetPeerSchoolData(ctx context.Context, nickName string) (model.PeerSchoolData, error) {
	var schoolData model.PeerSchoolData

	query, args, err := sq.Select("class_name", "parallel_name").
		From("participant").
		Where("login LIKE ?", "%"+nickName+"%").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return model.PeerSchoolData{}, fmt.Errorf("cannot configure query, err: %v", err)
	}

	err = r.conn.GetContext(ctx, &schoolData, query, args...)
	if err != nil {
		return model.PeerSchoolData{}, fmt.Errorf("cannot configure query, err: %v", err)
	}
	return schoolData, nil
}
