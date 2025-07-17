package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/s21platform/community-service/internal/model"
)

func (r *Repository) GetParticipantsLogin(ctx context.Context, limit, offset int64) ([]string, error) {
	var loginsParticipants []string

	query, args, err := sq.
		Select("login").
		From("participant").
		OrderBy("id ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to configure query, err: %v", err)
	}

	err = r.conn.SelectContext(ctx, &loginsParticipants, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants' logins, err: %v", err)
	}

	return loginsParticipants, nil
}

func (r *Repository) SetParticipantData(ctx context.Context, participantDataValue *model.ParticipantDataValue, login string) error {
	var campusID int
	err := r.conn.GetContext(ctx, &campusID, "SELECT id FROM campus WHERE campus_uuid = $1", participantDataValue.CampusUUID)
	if err != nil {
		return fmt.Errorf("failed to get campus_id by campus_uuid: %v", err)
	}

	query, args, err := sq.Update("participant").
		Set("campus_id", campusID).
		Set("class_name", participantDataValue.ClassName).
		Set("parallel_name", participantDataValue.ParallelName).
		Set("status", participantDataValue.Status).
		Set("exp_value", participantDataValue.ExpValue).
		Set("level", participantDataValue.Level).
		Set("exp_to_next_level", participantDataValue.ExpToNextLevel).
		Set("skills", participantDataValue.Skills).
		Set("crp", participantDataValue.PeerCodeReviewPoints).
		Set("prp", participantDataValue.PeerReviewPoints).
		Set("coins", participantDataValue.Coins).
		Set("badges", participantDataValue.Badges).
		//Set("tribe_id", tribeID).
		Where(sq.Eq{"login": login}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %v", err)
	}

	_, err = r.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %v", err)
	}

	return nil
}
