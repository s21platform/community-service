package main

import (
	"context"
	"sync"

	"github.com/s21platform/community-service/internal/service/campus"
)

func main() {
	// Создаём контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	// Завершаем работу воркера через контекст
	defer cancel()
	var wg sync.WaitGroup

	campusWorker := campus.New()

	wg.Add(1)

	// Запускаем горутину
	go func() {
		campusWorker.Run(ctx, &wg)
	}()

	// Даем время воркеру завершиться
	wg.Wait()
}