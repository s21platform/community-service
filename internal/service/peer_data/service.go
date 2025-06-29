package service

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

func (s *School) RunPeerWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("ParticipantDataWorker")

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("participant uploading worker shutting down")
			return

		case <-ticker.C:
			lastUpdate, err := s.rR.GetByKey(ctx, string(config.KeyParticipantDataLastUpdated))
			if err != nil {
				logger.Error(fmt.Sprintf("failed to get last update time, err: %v", err))
			}

			if lastUpdate == "" {
				err := s.uploadDataParticipant(ctx)
				if err != nil {
					logger.Error(fmt.Sprintf("failed to upload participants, err: %v", err))
				}

				err = s.rR.Set(ctx, string(config.KeyParticipantDataLastUpdated), "upd", time.Hour*24)
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
			return fmt.Errorf("failed to get participant logins, err: %v", err)
		}
		if len(logins) == 0 {
			break
		}

		for _, login := range logins {
			participantData, err := s.sC.GetParticipantData(ctx, login)
			if err != nil {
				logger.Error(fmt.Sprintf("failed to get participant data for login %s, err: %v", login, err))
				continue
			}

			if participantData == nil {
				mtx.Increment("update_paticipant_data.data_exists")
				continue
			}

			err = s.dbR.SaveParticipantData(ctx, participantData, login)
			if err != nil {
				logger.Error(fmt.Sprintf("failed to save participant data for login %s, err: %v", login, err))
			}
		}

		offset += limit
	}

	return nil

}
