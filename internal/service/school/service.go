package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/community-service/internal/config"
)

const (
	peerLimit = 1000
)

type School struct {
	sC  SchoolC
	dbR DbRepo
	rR  RedisRepo
}

func New(school SchoolC, dbR DbRepo, rR RedisRepo) *School {
	return &School{
		sC:  school,
		dbR: dbR,
		rR:  rR,
	}
}

func (s *School) RunPeerWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("RunPeerWorker")

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("participant worker done")

		case <-ticker.C:

			//case <-time.After(time.Hour * 24 * 30):
			campuses, err := s.dbR.GetCampusUuids(ctx)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot get campuses, err: %v", err))
				log.Fatalf("cannot get campuses, err: %v", err)
				return
			}

			var offset int64
			for _, campus := range campuses {
				offset = 0

				for {
					peerLogins, err := s.sC.GetPeersByCampusUuid(ctx, campus, peerLimit, offset)
					if err != nil {
						logger.Error(fmt.Sprintf("cannot get peer logins from school client, err: %v", err))
						log.Fatalf("cannot get peer logins from school client, err: %v", err)
						return
					}

					err = s.dbR.AddPeerLogins(ctx, peerLogins)
					if err != nil {
						logger.Error(fmt.Sprintf("cannot save peer logins, err: %v", err))
						log.Fatalf("cannot save peer logins, err: %v", err)
						return
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
