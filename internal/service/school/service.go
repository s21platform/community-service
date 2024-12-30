package service

import (
	"context"
	"github.com/s21platform/community-service/internal/rpc"
	"log"
	"sync"
	"time"
)

const (
	peerLimit = 1000
)

type School struct {
	school SchoolS
	dbR    rpc.DbRepo
	//нужно добавить редис
}

func New(school SchoolS, dbR rpc.DbRepo) *School {
	return &School{
		school: school,
		dbR:    dbR,
	}
}

func (s *School) RunPeerWorker(ctx context.Context, wg *sync.WaitGroup) {
	//откуда брять кампусы? из нашей бд или из школы тоже дергать ручки а здесь вызывать клиента?
	//get capmus uuids

	campuses := []string{"6bfe3c56-0211-4fe1-9e59-51616caac4dd"}

	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			//add logger
			log.Println("School service worker shutting down")

		case <-time.After(time.Hour * 24 * 30):
			for _, campus := range campuses {

				//как понять что дошли до максимума?
				var offset int64 = 0
				for offset < 10000 {
					peerLogins, err := s.school.GetPeersByCampusUuid(ctx, campus, peerLimit, offset)
					if err != nil {
						log.Fatalf("cannot get peer logins from school client, err: %v", err)
					}

					err = s.dbR.AddPeerLogins(ctx, peerLogins)
					if err != nil {
						log.Fatalf("cannot save peer logins, err: %v", err)
					}

					offset += peerLimit
				}

			}
		}
	}

}
