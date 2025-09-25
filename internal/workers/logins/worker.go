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

	metric := pkg.FromContext(ctx, config.KeyMetrics)

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger_lib.Info(ctx, "participant uploading worker shutting down")

		case <-ticker.C:
			lastUpdate, err := w.rR.GetByKey(ctx, config.KeyLoginsLastUpdated)
			if err != nil {
				logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get last update time")
			}

			if lastUpdate == "" {
				start := time.Now()
				err := w.uploadLogins(ctx)
				if err != nil {
					logger_lib.Error(logger_lib.WithError(ctx, err), "failed to upload logins")
				}
				metric.Duration(time.Since(start).Milliseconds(), "upload_peers.duration.worker")

				err = w.rR.Set(ctx, config.KeyLoginsLastUpdated, time.Now().Add(time.Hour*24*30).Format(time.RFC3339), time.Hour*24*15)
				if err != nil {
					logger_lib.Error(logger_lib.WithError(ctx, err), "failed to set last update time")
				}
			}
			logger_lib.Info(ctx, "participant worker done")
		}
	}
}

func (w *Worker) uploadLogins(ctx context.Context) error {
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
			time.Sleep(1 * time.Second)
			peerLogins, err := w.sC.GetPeersByCampusUuid(ctx, campus, peerLimit, offset)
			if err != nil {
				logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get chank of peer logins from school client")
				continue
			}

			for _, nickname := range peerLogins {
				login, err := w.dbR.GetPeerByLogin(ctx, nickname)
				if err != nil {
					metric.Increment("upload_peers.error")
					logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get peer by login")
					continue
				}
				if login.Nickname == "" {
					err = w.dbR.SetNickname(ctx, nickname)
					if err != nil {
						metric.Increment("upload_peers.error")
						logger_lib.Error(logger_lib.WithError(ctx, err), "failed to set peer by login")
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
			logger_lib.Info(ctx, "upload peer logins done")
			time.Sleep(2 * time.Second)
		}
		logger_lib.Info(ctx, fmt.Sprintf("read: %d peers (for campus: %s)", counter, campus))
	}
	return nil
}
