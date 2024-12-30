package service

import (
	"context"
	"log"
	"sync"
	"time"
)

const (
	peerLimit = 1000
)

type School struct {
	sC  SchoolC
	dbR DbRepo
	//add redis
}

func New(school SchoolC, dbR DbRepo) *School {
	return &School{
		sC:  school,
		dbR: dbR,
	}
}

func (s *School) RunPeerWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			//add logger
			log.Println("School service worker shutting down")

		case <-time.After(time.Hour * 24 * 30):
			campuses, err := s.dbR.GetCampusUuids(ctx)
			if err != nil {
				log.Fatalf("cannot get campuses, err: %v", err)
			}

			var offset int64
			for _, campus := range campuses {
				offset = 0

				for {
					peerLogins, err := s.sC.GetPeersByCampusUuid(ctx, campus, peerLimit, offset)
					if err != nil {
						log.Fatalf("cannot get peer logins from school client, err: %v", err)
					}

					err = s.dbR.AddPeerLogins(ctx, peerLogins)
					if err != nil {
						log.Fatalf("cannot save peer logins, err: %v", err)
					}

					if len(peerLogins) < peerLimit {
						break
					}

					offset += peerLimit
				}

			}
		}
	}

}
