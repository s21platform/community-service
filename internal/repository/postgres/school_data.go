package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/s21platform/community-service/internal/model"
)

func (r *Repository) GetPeerSchoolData(ctx context.Context, nickName string) (model.PeerSchoolData, error) {
	var schoolData model.PeerSchoolData
	login := fmt.Sprintf("%s@student.21-school.ru", nickName)

	query, args, err := sq.Select("class_name", "parallel_name").
		From("participant").
		Where(sq.Eq{"login": login}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return model.PeerSchoolData{}, fmt.Errorf("cannot configure query, err: %v", err)
	}

	err = r.conn.GetContext(ctx, &schoolData, query, args...)
	if err != nil {
		return model.PeerSchoolData{}, err
	}
	return schoolData, nil
}
