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

func (r *Repository) SetParticipantData(ctx context.Context, participantDataValue *model.ParticipantDataValue, login string) error {
	var campusID int
	err := r.conn.GetContext(ctx, &campusID, "SELECT id FROM campus WHERE campus_uuid = $1", participantDataValue.CampusUUID)
	if err != nil {
		return fmt.Errorf("failed to get campus_id by campus_uuid: %v", err)
	}

	query, args, err := sq.Insert("participant").
		Columns("login", "campus_id", "class_name", "parallel_name", "status", "exp_value", "level", "exp_to_next_level", "skills", "crp", "prp", "coins", "badges").
		Values(login, campusID, participantDataValue.ClassName, participantDataValue.ParallelName, participantDataValue.Status,
			participantDataValue.ExpValue, participantDataValue.Level, participantDataValue.ExpToNextLevel, participantDataValue.Skills,
			participantDataValue.PeerCodeReviewPoints, participantDataValue.PeerReviewPoints, participantDataValue.Coins, participantDataValue.Badges).
		Suffix(`
        ON CONFLICT (login) DO UPDATE SET
            campus_id = EXCLUDED.campus_id,
            class_name = EXCLUDED.class_name,
            parallel_name = EXCLUDED.parallel_name,
            status = EXCLUDED.status,
            exp_value = EXCLUDED.exp_value,
            level = EXCLUDED.level,
            exp_to_next_level = EXCLUDED.exp_to_next_level,
            skills = EXCLUDED.skills,
            crp = EXCLUDED.crp,
            prp = EXCLUDED.prp,
            coins = EXCLUDED.coins,
            badges = EXCLUDED.badges
    `).
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
