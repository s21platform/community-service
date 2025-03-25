package school

import (
	"context"
	"fmt"
	"log"

	school "github.com/s21platform/school-proto/school-proto"
	// "golang.org/x/tools/go/types/objectpath"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/s21platform/community-service/internal/config"
)

type Handle struct {
	client school.SchoolServiceClient
}

func (h *Handle) GetPeersByCampusUuid(ctx context.Context, campusUuid string, limit, offset int64) ([]string, error) {
	peers, err := h.client.GetPeers(ctx, &school.GetPeersIn{
		CampusUuid: campusUuid,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot get peers: %v", err)
	}

	return peers.Peer, err
}

func (h *Handle) GetParticipantData(ctx context.Context, peersLogin string) ( *school.GetParticipantDataOut , error) {
	peerLogin, err := h.client.GetParticipantData(ctx, &school.GetParticipantDataIn{
		Login: peersLogin,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot get peers: %v", err)
	}

	return peerLogin, err
}

func MustConnect(cfg *config.Config) *Handle {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.School.Host, cfg.School.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to school service: %v", err)
	}
	client := school.NewSchoolServiceClient(conn)
	return &Handle{client: client}
}

