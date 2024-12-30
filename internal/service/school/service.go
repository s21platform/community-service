package service

import (
	"context"
	"sync"
)

type School struct {
	school SchoolS
}

func New(school SchoolS) *School {
	return &School{school: school}
}

func (s *School) RunPeerWorker(ctx context.Context, wg *sync.WaitGroup) {
	//get capmus uuids

	defer wg.Done()

}
