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
		Select("nickname").
		From("login").
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

func (r *Repository) ParticipantData(ctx context.Context, login string) (*model.Participant, error) {
	var participant model.Participant

	query, args, err := sq.Select("login", "status").
		From("participant").
		Where(sq.Eq{"login": login}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build exists query: %v", err)
	}

	err = r.conn.GetContext(ctx, &participant, query, args...)
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

func (r *Repository) InsertParticipantData(ctx context.Context, participantDataValue *model.ParticipantDataValue, login string, campusID int64) error {
	query, args, err := sq.Insert("participant").
		Columns("login", "campus_id", "class_name", "parallel_name", "status", "exp_value", "level", "exp_to_next_level", "skills", "crp", "prp", "coins", "badges", "tribe_id").
		Values(login, campusID, participantDataValue.ClassName, participantDataValue.ParallelName, participantDataValue.Status,
			participantDataValue.ExpValue, participantDataValue.Level, participantDataValue.ExpToNextLevel, participantDataValue.Skills,
			participantDataValue.PeerCodeReviewPoints, participantDataValue.PeerReviewPoints, participantDataValue.Coins, participantDataValue.Badges, participantDataValue.TribeID).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %v", err)
	}

	_, err = r.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert participant: %v", err)
	}

	return nil
}

func (r *Repository) UpdateParticipantData(ctx context.Context, participantDataValue *model.ParticipantDataValue, login string, campusID int64) error {
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
		Set("tribe_id", participantDataValue.TribeID).
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
