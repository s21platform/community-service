package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/community-service/internal/config"
)

const (
	peerLimit = 1000
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
	logger.AddFuncName("RunPeerWorker")

	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("participant worker uploading done")

		case <-ticker.C:
			timeCheck, err := s.updateTimeCheck(ctx)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot check time since last update, err: %v", err))
			}

			if timeCheck {
				err := s.uploadParticipants(ctx)
				if err != nil {
					logger.Error(fmt.Sprintf("cannot upload participants, err: %v", err))
				}

				err = s.setUpdateTime(ctx)
				if err != nil {
					logger.Error(fmt.Sprintf("cannot save participant last updated, err: %v", err))
				}
			}
		}
	}
}

func (s *School) setUpdateTime(ctx context.Context) error {
	timeUpdated := time.Now().Format(time.RFC3339)
	err := s.rR.Set(ctx, config.KeyParticipantLastUpdated, timeUpdated, time.Hour*24*60)
	if err != nil {
		return fmt.Errorf("cannot save participant last updated, err: %v", err)
	}
	return nil
}

// если прошел месяц после последнего обновления - возвращает тру или если записей еще не было
func (s *School) updateTimeCheck(ctx context.Context) (bool, error) {
	lastUpdate, err := s.rR.GetByKey(ctx, config.KeyParticipantLastUpdated)
	if err != nil {
		return false, fmt.Errorf("cannot get last update time, err: %v", err)
	}
	if lastUpdate == "" {
		return true, nil //если последнее время еще не добавлено
	}

	lastUpdateTime, err := time.Parse(time.RFC3339, lastUpdate)
	if err != nil {
		return false, fmt.Errorf("cannot parse time, err: %v", err)
	}

	if time.Now().After(lastUpdateTime.AddDate(0, 1, 0)) {
		return true, nil
	}
	return false, nil
}

func (s *School) uploadParticipants(ctx context.Context) error {
	campuses, err := s.dbR.GetCampusUuids(ctx)
	if err != nil {
		return fmt.Errorf("cannot get campuses, err: %v", err)
	}

	var offset int64
	for _, campus := range campuses {
		offset = 0

		for {
			peerLogins, err := s.sC.GetPeersByCampusUuid(ctx, campus, peerLimit, offset)
			if err != nil {
				return fmt.Errorf("cannot get peer logins from school client, err: %v", err)
			}

			err = s.dbR.AddPeerLogins(ctx, peerLogins)
			if err != nil {
				return fmt.Errorf("cannot save peer logins, err: %v", err)
			}

			if len(peerLogins) < peerLimit {
				break
			}
			offset += peerLimit
		}
	}
	return nil
}
