package campus

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/s21platform/community-service/internal/config"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"
)

type Worker struct {
	sC  SchoolClient
	dbR DbRepo
	rR  RedisRepo
}

func New(sC SchoolClient, dbR DbRepo, rR RedisRepo) *Worker {
	return &Worker{
		sC:  sC,
		dbR: dbR,
		rR:  rR,
	}
}

func (w *Worker) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("WorkerCampusRun")

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("campus uploading worker shutting down")

		case <-ticker.C:
			lastUpdate, err := w.rR.GetByKey(ctx, config.KeyCampusesLastUpdated)
			if err != nil {
				logger.Error(fmt.Sprintf("failed to get last update time, err: %v", err))
				continue
			}
			if lastUpdate != "" {
				continue
			}

			err = w.process(ctx)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot upload campuses, err: %v", err))
				continue
			}

			err = w.rR.Set(ctx, string(config.KeyCampusesLastUpdated), "upd", time.Hour*5)
			if err != nil {
				logger.Error(fmt.Sprintf("failed to save campuses last updated, err: %v", err))
				continue
			}

			logger.Info("campuses worker done")
		}
	}
}

func (w *Worker) process(ctx context.Context) error {
	mtx := pkg.FromContext(ctx, config.KeyMetrics)
	campuses, err := w.sC.GetCampuses(ctx)
	if err != nil {
		return fmt.Errorf("failed to get campuses, err: %v", err)
	}

	for _, campus := range campuses {
		cpm, err := w.dbR.GetCampusByUUID(ctx, campus.Uuid)
		if err != nil {
			return fmt.Errorf("failed to check campus exist, err: %v", err)
		}

		if cpm != nil {
			mtx.Increment("campus_exists")
			continue
		}

		err = w.dbR.SetCampus(ctx, campus)
		if err != nil {
			return fmt.Errorf("failed to create campus, err: %v", err)
		}
		mtx.Increment("campus_created")
	}

	return nil
}
