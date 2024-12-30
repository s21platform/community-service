package main

import (
	"context"
	"github.com/s21platform/community-service/internal/clients/school"
	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/service/school"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.MustLoad()
	schoolClient := school.MustConnect(cfg)
	peerWorker := service.New(schoolClient)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		peerWorker.RunPeerWorker(ctx, wg)
	}()
	wg.Wait()
}
