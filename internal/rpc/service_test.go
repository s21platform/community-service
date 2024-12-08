package rpc_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	community_proto "github.com/s21platform/community-proto/community-proto"
	"github.com/s21platform/community-service/internal/model"
	"github.com/s21platform/community-service/internal/rpc"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestServer_GetPeerSchoolData(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepo := rpc.NewMockDbRepo(controller)

	t.Run("get_peer_school_data_ok", func(t *testing.T) {
		expectedData := model.PeerSchoolData{ClassName: "test-class", ParallelName: "test-parallel"}
		nickName := "aboba"
		mockRepo.EXPECT().GetPeerSchoolData(gomock.Any(), nickName).Return(expectedData, nil)

		s := rpc.New(mockRepo)
		data, err := s.GetPeerSchoolData(ctx, &community_proto.GetSchoolDataIn{NickName: nickName})
		assert.NoError(t, err)
		assert.Equal(t, data, &community_proto.GetSchoolDataOut{ClassName: expectedData.ClassName, ParallelName: expectedData.ParallelName})
	})

	t.Run("get_peer_school_data_err", func(t *testing.T) {
		nickName := "aboba"
		expectedErr := errors.New("select err")
		mockRepo.EXPECT().GetPeerSchoolData(gomock.Any(), nickName).Return(model.PeerSchoolData{}, expectedErr)

		s := rpc.New(mockRepo)

		data, err := s.GetPeerSchoolData(ctx, &community_proto.GetSchoolDataIn{NickName: nickName})
		assert.Nil(t, data)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "select err")
	})
}

func TestServer_IsPeerExist(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepo := rpc.NewMockDbRepo(controller)

	t.Run("get_peer_status_ok", func(t *testing.T) {
		expectedStatus := "ACTIVE"
		email := "aboba@student.21-school.ru"
		mockRepo.EXPECT().GetPeerStatus(gomock.Any(), email).Return(expectedStatus, nil)

		s := rpc.New(mockRepo)
		isExist, err := s.IsPeerExist(ctx, &community_proto.EmailIn{Email: email})
		assert.NoError(t, err)
		assert.True(t, isExist.IsExist)
	})

	t.Run("get_peer_status_not_found", func(t *testing.T) {
		expectedStatus := "NOT ACTIVE"
		email := "aboba@student.21-school.ru"
		mockRepo.EXPECT().GetPeerStatus(gomock.Any(), email).Return(expectedStatus, nil)

		s := rpc.New(mockRepo)
		isExist, err := s.IsPeerExist(ctx, &community_proto.EmailIn{Email: email})
		assert.NoError(t, err)
		assert.False(t, isExist.IsExist)
	})

	t.Run("get_peer_status_err", func(t *testing.T) {
		email := "aboba@student.21-school.ru"
		expectedErr := errors.New("select err")
		mockRepo.EXPECT().GetPeerStatus(gomock.Any(), email).Return("", expectedErr)

		s := rpc.New(mockRepo)
		isExist, err := s.IsPeerExist(ctx, &community_proto.EmailIn{Email: email})

		assert.Nil(t, isExist)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "select err")
	})
}
