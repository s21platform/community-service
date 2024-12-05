package campus

import (
	"context"
	"log"
	"sync"
	"time"
)

type Campus struct {
	
}

func New() *Campus {
	return &Campus{}
}

func (s *Campus) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()	
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