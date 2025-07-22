package main

import (
	"context"
	"log"
	"sync"

	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/community-service/internal/client/school"
	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/repository/postgres"
	"github.com/s21platform/community-service/internal/repository/redis"
	"github.com/s21platform/community-service/internal/workers/campus"
)

func main() {
	cfg := config.MustLoad()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)
	ctx = context.WithValue(ctx, config.KeyLogger, logger)

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, cfg.Service.Name, cfg.Platform.Env)
	if err != nil {
		log.Fatalf("failed to create metrics: %s", err)
	}
	defer metrics.Disconnect()
	ctx = context.WithValue(ctx, config.KeyMetrics, metrics)

	schoolClient := school.MustConnect(cfg)
	dbRepo := postgres.New(cfg)
	defer dbRepo.Close()
	redisRepo := redis.New(cfg)
	campusWorker := campus.New(schoolClient, dbRepo, redisRepo)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		campusWorker.Run(ctx, wg)
	}()
	wg.Wait()
}
