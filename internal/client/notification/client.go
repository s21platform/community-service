package notification

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/notification-service/pkg/notification"
)

type Client struct {
	client notification.NotificationServiceClient
}

func New(cfg *config.Config) *Client {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.Notification.Host, cfg.Notification.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to notification client: %v", err)
	}
	client := notification.NewNotificationServiceClient(conn)
	return &Client{client: client}
}

func (c *Client) SendEduCode(ctx context.Context, email, code string) error {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("uuid", ctx.Value(config.KeyUUID).(string)))

	_, err := c.client.SendEduCode(ctx, &notification.SendEduCodeIn{
		Email: email,
		Code:  code,
	})
	if err != nil {
		return fmt.Errorf("failed to call notification service: %v", err)
	}

	return nil
}
