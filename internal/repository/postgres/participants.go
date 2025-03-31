package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	// "database/sql"
	// "errors"
	"fmt"
	// `github.com/s21platform/community-service/internal/model`
	sq "github.com/Masterminds/squirrel"
	school "github.com/s21platform/school-proto/school-proto"
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


func (r *Repository) SetParticipantData(ctx context.Context, participantData *school.GetParticipantDataOut, login string) error {
	tx, err := r.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("cannot start transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	skillsJSON, err := json.Marshal(participantData.Skills)
	if err != nil {
		return fmt.Errorf("cannot marshal skills: %v", err)
	}

	badgesJSON, err := json.Marshal(participantData.Badges)
	if err != nil {
		return fmt.Errorf("cannot marshal badges: %v", err)
	}

	queryBase := sq.Insert("participant").
		Columns(
			"login", "campus_id", "class_name", "parallel_name",
			"status", "exp_value", "level", "exp_to_next_level", "skills",
			"crp", "prp", "coins", "badges").
		Values(
			login, 
			participantData.CampusUuid,
			participantData.ClassName, 
			participantData.ParallelName,
			participantData.Status, 
			participantData.ExpValue, 
			participantData.Level,
			participantData.ExpToNextLevel, 
			skillsJSON,
			participantData.PeerCodeReviewPoints,
			participantData.PeerReviewPoints, 
			participantData.Coins, 
			badgesJSON).
		Suffix(`
			ON CONFLICT (login)
			DO UPDATE SET
				exp_value = EXCLUDED.exp_value,
				level = EXCLUDED.level,
				exp_to_next_level = EXCLUDED.exp_to_next_level,
				skills = EXCLUDED.skills,
				crp = EXCLUDED.crp,
				prp = EXCLUDED.prp,
				coins = EXCLUDED.coins,
				badges = EXCLUDED.badges`).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBase.ToSql()
	if err != nil {
		return fmt.Errorf("cannot configure query: %v", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("cannot execute query: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cannot commit transaction: %v", err)
	}

	return nil
}