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

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("participant uploading worker shutting down")
			return

		case <-ticker.C:
			lastUpdate, err := s.rR.GetByKey(ctx, config.KeyParticipantLastUpdated)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot get last update time, err: %v", err))
			}

			if lastUpdate == "" {
				err := s.uploadParticipants(ctx) // заменить на новую функцию
				if err != nil {
					logger.Error(fmt.Sprintf("cannot upload participants, err: %v", err))
				}

				err = s.rR.Set(ctx, config.KeyParticipantLastUpdated, "upd", time.Hour*24*30)
				if err != nil {
					logger.Error(fmt.Sprintf("cannot save participant last updated, err: %v", err))
				}
			}
			logger.Info("participant worker done")
		}
	}
}

func (s *School) uploadInformationAboutParticipant(ctx context.Context) error {
	campuses, err := s.dbR.GetCampusUuids(ctx)
	if err != nil {
		return fmt.Errorf("cannot get campuses, err: %v", err)
	}

	for _, campus := range campuses {// по логинам
		offset := int64(0)

		for {
			peerLogins, err := s.sC.GetParticipantData(ctx, )
			if err != nil {
				return fmt.Errorf("cannot get peer logins from school client, err: %v", err)
			}

			err = s.dbR.AddPeerLogins(ctx, peerLogins)
			if err != nil {
				return fmt.Errorf("cannot save peer logins, err: %v", err)
			}

			for _, login := range peerLogins {
				participantData, err := s.sC.GetParticipantData(ctx, login)
				if err != nil {
					logger_lib.FromContext(ctx, config.KeyLogger).Error(fmt.Sprintf("failed to fetch participant data for %s: %v", login, err))
					continue
				}

				err = s.dbR.SaveParticipantData(ctx, participantData)
				if err != nil {
					logger_lib.FromContext(ctx, config.KeyLogger).Error(fmt.Sprintf("failed to save participant data for %s: %v", login, err))
				}
			}

			if len(peerLogins) < peerLimit {
				break
			}
			offset += peerLimit
		}
	}
	return nil
}