package service

import (
	"context"
	"github.com/s21platform/community-service/internal/rpc"
	"log"
	"sync"
	"time"
)

type School struct {
	school SchoolS
	dbR rpc.DbRepo
}

func New(school SchoolS, dbR rpc.DbRepo) *School {
	return &School{
		school: school,
		dbR: dbR,
	}
}

func (s *School) RunPeerWorker(ctx context.Context, wg *sync.WaitGroup) {
	//get capmus uuids

	campuses := []string{"6bfe3c56-0211-4fe1-9e59-51616caac4dd"}

	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			//add logger
			log.Println("School service worker shutting down")

		case <- time.After(time.Hour * 24 * 30):
			s.dbR.
		}
	}

}
