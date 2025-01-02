package main

import (
	"context"
	logger_lib "github.com/s21platform/logger-lib"
	"sync"

	"github.com/s21platform/community-service/internal/clients/school"
	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/repository/postgres"
	service "github.com/s21platform/community-service/internal/service/school"
)

func main() {
	cfg := config.MustLoad()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)
	ctx = context.WithValue(ctx, config.KeyLogger, logger)

	schoolClient := school.MustConnect(cfg)
	dbRepo := postgres.New(cfg)
	peerWorker := service.New(schoolClient, dbRepo)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		peerWorker.RunPeerWorker(ctx, wg)
	}()
	wg.Wait()
}
