package main

import (
	"context"
	"log"
	"time"
)

func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Worker is stopping...")
			return 
		case <-time.After(30 * 24 * time.Hour): 
			log.Println("work")
		}
	}
}

func main() {
	// Создаём контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())

	// Запускаем горутину
	go worker(ctx)

	// Завершаем работу воркера через контекст
	cancel()

	// Даем время воркеру завершиться
	time.Sleep(1 * time.Second)
}