package main

import (
	"context"
	"sync"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/community-service/internal/clients/school"
	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/repository/postgres"
	service "github.com/s21platform/community-service/internal/service/school"
)

func main() {
	cfg := config.MustLoad()
	logger := logger_lib.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
