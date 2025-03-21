package school

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	school "github.com/s21platform/school-proto/school-proto"
	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/model"
)

type Client struct {
	client school.SchoolServiceClient
}

func (c *Client) GetPeersByCampusUuid(ctx context.Context, campusUuid string, limit, offset int64) ([]string, error) {
	peers, err := c.client.GetPeers(ctx, &school.GetPeersIn{
		CampusUuid: campusUuid,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot get peers: %v", err)
	}

	return peers.Peer, err
}

func (c *Client) GetCampuses(ctx context.Context) ([]model.Campus, error) {
	campuses, err := c.client.GetCampuses(ctx, &school.Empty{})
	if err != nil {
		return nil, fmt.Errorf("cannot get campuses: %v", err)
	}
	var result []model.Campus
	for _, campus := range campuses.Campuses {
		result = append(result, model.Campus{
			Uuid:      campus.CampusUuid,
			FullName:  campus.FullName,
			ShortName: campus.ShortName,
		})
	}
	return result, nil
}

func MustConnect(cfg *config.Config) *Client {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.School.Host, cfg.School.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to school service: %v", err)
	}
	client := school.NewSchoolServiceClient(conn)
	return &Client{client: client}
}
