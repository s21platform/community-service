package peerdata

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/model"
)

const (
	limit = 10000
)

type School struct {
	sC  SchoolC
	dbR DbRepo
	rR  RedisRepo
}

func New(school SchoolC, dbR DbRepo, rR RedisRepo) *School {
	return &School{
		sC:  school,
		dbR: dbR,
		rR:  rR,
	}
}

func (s *School) RunParticipantWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("ParticipantDataWorker")
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("participant uploading worker shutting down")
			return

		case <-ticker.C:
			lastUpdate, err := s.rR.GetByKey(ctx, config.KeyParticipantDataLastUpdated)
			if err != nil {
				logger.Error(fmt.Sprintf("failed to get last update time, err: %v", err))
			}

			if lastUpdate == "" {
				err := s.uploadDataParticipant(ctx)
				if err != nil {
					logger.Error(fmt.Sprintf("failed to upload participants, err: %v", err))
				}

				err = s.rR.Set(ctx, config.KeyParticipantDataLastUpdated, "upd", 5*time.Hour)
				if err != nil {
					logger.Error(fmt.Sprintf("failed to save participant last updated, err: %v", err))
				}
			}
			logger.Info("participant worker done")
		}
	}
}

func (s *School) uploadDataParticipant(ctx context.Context) error {
	var offset int64
	mtx := pkg.FromContext(ctx, config.KeyMetrics)
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	for {
		logins, err := s.dbR.GetParticipantsLogin(ctx, limit, offset)
		if err != nil {
			mtx.Increment("update_participant_data.edu_error")
			return fmt.Errorf("failed to get participant logins, err: %v", err)
		}
		if len(logins) == 0 {
			mtx.Increment("update_participant_data.empty_login_list")
			break
		}

		for _, login := range logins {
			exists := true
			participant, err := s.dbR.ParticipantData(ctx, login)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					mtx.Increment("update_participant_data.error_get_participant")
					logger.Error(fmt.Sprintf("failed to check participant existance: %v", err))
					continue
				}
				exists = false
			}
			if participant != nil && (participant.Status == model.ParticipantStatusBlocked || participant.Status == model.ParticipantStatusFrozen || participant.Status == model.ParticipantStatusExpelled) {
				mtx.Increment("update_participant_data.skip_not_active")
				continue
			}
			time.Sleep(200 * time.Millisecond)
			participantData, err := s.sC.GetParticipantData(ctx, login)
			if err != nil {
				if strings.Contains(err.Error(), "Invalid token") {
					mtx.Increment("update_participant_data.invalid_token")
				} else if strings.Contains(err.Error(), "Too many requests") {
					mtx.Increment("update_participant_data.too_many_requests")
				} else {
					mtx.Increment("update_participant_data.unknown_error")
				}

				logger.Error(fmt.Sprintf("failed to get participant data for login %s, err: %v", login, err))
				continue
			}

			if participantData == nil {
				mtx.Increment("update_participant_data.not_exists")
				continue
			}

			campus, err := s.dbR.GetCampusByUUID(ctx, participantData.CampusUUID)
			if err != nil {
				mtx.Increment("update_participant_data.error_get_campus")
				logger.Error(fmt.Sprintf("failed to get participant campus: %v", err))
				continue
			}
			participantData.TribeID = 1
			if !exists {
				err = s.dbR.InsertParticipantData(ctx, participantData, login, campus.Id)
			} else {
				err = s.dbR.UpdateParticipantData(ctx, participantData, login, campus.Id)
			}
			if err != nil {
				mtx.Increment("update_participant_data.not_save")
				logger.Error(fmt.Sprintf("failed to save participant data for login %s, err: %v", login, err))
				continue
			}
			mtx.Increment("update_participant_data.ok")
		}

		offset += limit
	}
	mtx.Increment("update_participant_data.finish_upload")
	return nil
}
