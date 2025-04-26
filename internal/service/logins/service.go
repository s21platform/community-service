package logins

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

type Worker struct {
	sC  SchoolClient
	dbR DbRepo
	rR  RedisRepo
}

func New(school SchoolClient, dbR DbRepo, rR RedisRepo) *Worker {
	return &Worker{
		sC:  school,
		dbR: dbR,
		rR:  rR,
	}
}

func (w *Worker) RunPeerWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("RunPeerWorker")

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("participant uploading worker shutting down")

		case <-ticker.C:
			lastUpdate, err := w.rR.GetByKey(ctx, config.KeyLoginsLastUpdated)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot get last update time, err: %v", err))
			}

			if lastUpdate == "" {
				err := w.UploadParticipants(ctx)
				if err != nil {
					logger.Error(fmt.Sprintf("cannot upload participants, err: %v", err))
				}

				err = w.rR.Set(ctx, config.KeyLoginsLastUpdated, time.Now().Add(time.Hour*24*30).Format(time.RFC3339), time.Hour*24*30)
				if err != nil {
					logger.Error(fmt.Sprintf("cannot save participant last updated, err: %v", err))
				}
			}
			logger.Info("participant worker done")
		}
	}
}

func (w *Worker) UploadParticipants(ctx context.Context) error {
	campuses, err := w.dbR.GetCampusUuids(ctx)
	if err != nil {
		return fmt.Errorf("cannot get campuses, err: %v", err)
	}

	var offset int64
	for _, campus := range campuses {
		offset = 0

		for {
			peerLogins, err := w.sC.GetPeersByCampusUuid(ctx, campus, peerLimit, offset)
			if err != nil {
				return fmt.Errorf("cannot get peer logins from logins client, err: %v", err)
			}

			err = w.dbR.AddPeerLogins(ctx, peerLogins)
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
