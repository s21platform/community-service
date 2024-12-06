package rpc_test

import (
	"github.com/golang/mock/gomock"
	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/rpc"
	"golang.org/x/net/context"
	"testing"
)

func TestServer_GetPeerSchoolData(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	uuid := "test_uuid"
	ctx = context.WithValue(ctx, config.KeyUUID, uuid)

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepo := rpc.NewMockDbRepo(controller)
}
