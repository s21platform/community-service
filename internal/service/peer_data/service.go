package peerdata

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/s21platform/community-service/internal/config"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"
)

const (
	limit = 1000
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

				err = s.rR.Set(ctx, config.KeyParticipantDataLastUpdated, "upd", time.Hour)
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
			participantData, err := s.sC.GetParticipantData(ctx, login)
			if err != nil {
				logger.Error(fmt.Sprintf("failed to get participant data for login %s, err: %v", login, err))
				continue
			}

			if participantData == nil {
				mtx.Increment("update_participant_data.not_exists")
				continue
			}

			exists, err := s.dbR.IsParticipantDataExists(ctx, login)
			if err != nil {
				mtx.Increment("update_participant_data.error_get_participant")
				logger.Error(fmt.Sprintf("failed to check participant existance: %v", err))
				continue
			}

			campus, err := s.dbR.GetCampusByUUID(ctx, participantData.CampusUUID)
			if err != nil {
				mtx.Increment("update_participant_data.error_get_campus")
				logger.Error(fmt.Sprintf("failed to get participant campus: %v", err))
				continue
			}
			participantData.TribeID = "no tribe yet"
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
			time.Sleep(10 * time.Millisecond)
		}

		offset += limit
	}

	return nil
}
