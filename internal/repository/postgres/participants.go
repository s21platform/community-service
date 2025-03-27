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

func (r *Repository) GetParticipantsLogin(ctx context.Context) ([]string, error) {
	var loginsParticipants []string
	query, args, err := sq.Select("login").
		From("participant").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot configure query, err: %v", err)
	}

	err = r.conn.SelectContext(ctx, &loginsParticipants, query, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot get participants' logins, err: %v", err)
	}

	return loginsParticipants, nil
}

func (r *Repository) SetParticipantData(ctx context.Context, participantData *school.GetParticipantDataOut) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot start transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	queryBase := sq.Insert("participants").
		Columns(
			"class_name", "parallel_name", "exp_value", "level", "exp_to_next_level",
			"campus_uuid", "status", "peer_review_points", "peer_code_review_points", "coins").
		Values(
			participantData.ClassName, participantData.ParallelName, participantData.ExpValue, participantData.Level,
			participantData.ExpToNextLevel, participantData.CampusUuid, participantData.Status,
			participantData.PeerReviewPoints, participantData.PeerCodeReviewPoints, participantData.Coins).
		Suffix("ON CONFLICT (campus_uuid) DO UPDATE SET exp_value = EXCLUDED.exp_value, level = EXCLUDED.level").
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBase.ToSql()
	if err != nil {
		return fmt.Errorf("cannot configure query: %v", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("cannot execute query: %v", err)
	}

	for _, skill := range participantData.Skills {
		skillQuery := sq.Insert("skills").
			Columns("campus_uuid", "skill_name", "skill_points").
			Values(participantData.CampusUuid, skill.Name, skill.Points).
			Suffix("ON CONFLICT (campus_uuid, skill_name) DO UPDATE SET skill_points = EXCLUDED.skill_points").
			PlaceholderFormat(sq.Dollar)

		skillSQL, skillArgs, err := skillQuery.ToSql()
		if err != nil {
			return fmt.Errorf("cannot configure skills query: %v", err)
		}

		_, err = tx.ExecContext(ctx, skillSQL, skillArgs...)
		if err != nil {
			return fmt.Errorf("cannot execute skills query: %v", err)
		}
	}


	for _, badge := range participantData.Badges {
		badgeQuery := sq.Insert("badges").
			Columns("campus_uuid", "badge_name", "receipt_date_time", "icon_url").
			Values(participantData.CampusUuid, badge.Name, badge.ReceiptDateTime, badge.IconURL).
			Suffix("ON CONFLICT (campus_uuid, badge_name) DO UPDATE SET receipt_date_time = EXCLUDED.receipt_date_time, icon_url = EXCLUDED.icon_url").
			PlaceholderFormat(sq.Dollar)

		badgeSQL, badgeArgs, err := badgeQuery.ToSql()
		if err != nil {
			return fmt.Errorf("cannot configure badges query: %v", err)
		}

		_, err = tx.ExecContext(ctx, badgeSQL, badgeArgs...)
		if err != nil {
			return fmt.Errorf("cannot execute badges query: %v", err)
		}
	}

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cannot commit transaction: %v", err)
	}

	return nil
}