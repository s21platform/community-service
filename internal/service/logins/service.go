package logins

import (
	"context"
	"fmt"
	"sync"
	"time"

	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"

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

	metric := pkg.FromContext(ctx, config.KeyMetrics)

	ticker := time.NewTicker(10 * time.Minute)
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
				logger.Info(fmt.Sprintf("Update Logins is start at %s", time.Now().Format(time.RFC3339)))
				start := time.Now()
				err := w.uploadLogins(ctx)
				if err != nil {
					logger.Error(fmt.Sprintf("cannot upload participants, err: %v", err))
				}
				metric.Duration(time.Since(start).Milliseconds(), "upload_peers.duration.worker")

				err = w.rR.Set(ctx, config.KeyLoginsLastUpdated, time.Now().Add(time.Hour*24*30).Format(time.RFC3339), time.Hour*24*30)
				if err != nil {
					logger.Error(fmt.Sprintf("cannot save participant last updated, err: %v", err))
				}
			}
			logger.Info("participant worker done")
		}
	}
}

func (w *Worker) uploadLogins(ctx context.Context) error {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	metric := pkg.FromContext(ctx, config.KeyMetrics)

	campuses, err := w.dbR.GetCampusUuids(ctx)
	if err != nil {
		return fmt.Errorf("cannot get campuses, err: %v", err)
	}

	var offset int64
	for _, campus := range campuses {
		offset = 0
		counter := 0

		for {
			peerLogins, err := w.sC.GetPeersByCampusUuid(ctx, campus, peerLimit, offset)
			if err != nil {
				logger.Error(fmt.Sprintf("failed to get chank of peer logins from school client - skip: %v", err))
				continue
			}

			for _, nickname := range peerLogins {
				login, err := w.dbR.GetPeerByLogin(ctx, nickname)
				if err != nil {
					metric.Increment("upload_peers.error")
					logger.Error(fmt.Sprintf("failed to get nickname from login table: %v", err))
					continue
				}
				if login.Nickname == "" {
					err = w.dbR.SetNickname(ctx, nickname)
					if err != nil {
						metric.Increment("upload_peers.error")
						logger.Error(fmt.Sprintf("failed to save peer nickname: %v", err))
						continue
					}
					metric.Increment("upload_peers.new")
					continue
				}
				metric.Increment("upload_peers.already_exist")
			}

			if len(peerLogins) < peerLimit {
				counter += len(peerLogins)
				break
			}
			offset += peerLimit
			counter += peerLimit
			logger.Info(fmt.Sprintf("iteration complete: campus: %s", campus))
			time.Sleep(3 * time.Second)
		}
		logger.Info(fmt.Sprintf("read: %d peers (for campus: %s)", counter, campus))
	}
	return nil
}
