package main

import (
	"context"
	"sync"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/community-service/internal/client/school"
	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/repository/postgres"
	"github.com/s21platform/community-service/internal/repository/redis"
	"github.com/s21platform/community-service/internal/service/logins"
)

func main() {
	cfg := config.MustLoad()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)
	ctx = context.WithValue(ctx, config.KeyLogger, logger)

	schoolClient := school.MustConnect(cfg)
	dbRepo := postgres.New(cfg)
	redisRepo := redis.New(cfg)
	peerWorker := logins.New(schoolClient, dbRepo, redisRepo)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		peerWorker.RunPeerWorker(ctx, wg)
	}()
	wg.Wait()
}
