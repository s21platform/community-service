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

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger_lib.Info(ctx, "campus uploading worker shutting down")

		case <-ticker.C:
			lastUpdate, err := w.rR.GetByKey(ctx, config.KeyCampusesLastUpdated)
			if err != nil {
				logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get last update time")
				continue
			}
			if lastUpdate != "" {
				continue
			}

			err = w.process(ctx)
			if err != nil {
				logger_lib.Error(logger_lib.WithError(ctx, err), "failed to process campuses")
				continue
			}

			err = w.rR.Set(ctx, config.KeyCampusesLastUpdated, "upd", time.Hour*24*15)
			if err != nil {
				logger_lib.Error(logger_lib.WithError(ctx, err), "failed to set last update time")
				continue
			}

			logger_lib.Info(ctx, "campuses worker done")
		}
	}
}

func (w *Worker) process(ctx context.Context) error {
	mtx := pkg.FromContext(ctx, config.KeyMetrics)
	campuses, err := w.sC.GetCampuses(ctx)
	if err != nil {
		mtx.Increment("upload_campus.failed_to_get")
		return fmt.Errorf("failed to get campuses, err: %v", err)
	}

	for _, campus := range campuses {
		cpm, err := w.dbR.GetCampusByUUID(ctx, campus.Uuid)
		if err != nil {
			mtx.Increment("upload_campus.failed_to_get_from_db")
			return fmt.Errorf("failed to check campus exist, err: %v", err)
		}

		if cpm != nil {
			mtx.Increment("upload_campus.already_exist")
			continue
		}

		err = w.dbR.SetCampus(ctx, campus)
		if err != nil {
			return fmt.Errorf("failed to create campus, err: %v", err)
		}
		mtx.Increment("upload_campus.new")
	}

	return nil
}
