package rpc_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	community_proto "github.com/s21platform/community-proto/community-proto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/s21platform/community-service/internal/model"
	"github.com/s21platform/community-service/internal/rpc"
)

var env = "prod"

func TestService_IsPeerExist(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepo := rpc.NewMockDbRepo(controller)

	t.Run("is_peer_exist_ok", func(t *testing.T) {
		login := "staff_login"
		var id int64 = 1

		mockRepo.EXPECT().GetStaffId(gomock.Any(), login).Return(id, nil)

		s := rpc.New(mockRepo, "prod")

		data, err := s.IsUserStaff(ctx, &community_proto.LoginIn{Login: login})
		assert.NoError(t, err)
		assert.Equal(t, data, &community_proto.IsUserStaffOut{IsStaff: true})
	})
}

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

		s := rpc.New(mockRepo, env)
		data, err := s.GetPeerSchoolData(ctx, &community_proto.GetSchoolDataIn{NickName: nickName})
		assert.NoError(t, err)
		assert.Equal(t, data, &community_proto.GetSchoolDataOut{ClassName: expectedData.ClassName, ParallelName: expectedData.ParallelName})
	})

	t.Run("get_peer_school_data_err", func(t *testing.T) {
		nickName := "aboba"
		expectedErr := errors.New("select err")
		mockRepo.EXPECT().GetPeerSchoolData(gomock.Any(), nickName).Return(model.PeerSchoolData{}, expectedErr)

		s := rpc.New(mockRepo, env)

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

		s := rpc.New(mockRepo, env)
		isExist, err := s.IsPeerExist(ctx, &community_proto.EmailIn{Email: email})
		assert.NoError(t, err)
		assert.True(t, isExist.IsExist)
	})

	t.Run("get_peer_status_not_found", func(t *testing.T) {
		expectedStatus := "NOT ACTIVE"
		email := "aboba@student.21-school.ru"
		mockRepo.EXPECT().GetPeerStatus(gomock.Any(), email).Return(expectedStatus, nil)

		s := rpc.New(mockRepo, env)
		isExist, err := s.IsPeerExist(ctx, &community_proto.EmailIn{Email: email})
		assert.NoError(t, err)
		assert.False(t, isExist.IsExist)
	})

	t.Run("get_peer_status_err", func(t *testing.T) {
		email := "aboba@student.21-school.ru"
		expectedErr := errors.New("select err")
		mockRepo.EXPECT().GetPeerStatus(gomock.Any(), email).Return("", expectedErr)

		s := rpc.New(mockRepo, env)
		isExist, err := s.IsPeerExist(ctx, &community_proto.EmailIn{Email: email})

		assert.Nil(t, isExist)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "select err")
	})
}

// limit и offset сейчас не используются
func TestServer_SearchPeers(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepo := rpc.NewMockDbRepo(controller)

	t.Run("search_peers_ok", func(t *testing.T) {
		expectedData := []*community_proto.SearchPeer{
			{Login: "aboba"},
			{Login: "abobaoba"},
			{Login: "aboo"},
		}
		substring := "ab"
		mockRepo.EXPECT().SearchPeersBySubstring(gomock.Any(), substring).Return(expectedData, nil)

		s := rpc.New(mockRepo, env)
		data, err := s.SearchPeers(ctx, &community_proto.SearchPeersIn{Substring: substring})
		assert.NoError(t, err)
		assert.Equal(t, data, &community_proto.SearchPeersOut{SearchPeers: expectedData})
	})

	t.Run("search_peers_err", func(t *testing.T) {
		expectedErr := errors.New("select err")
		substring := "ab"
		mockRepo.EXPECT().SearchPeersBySubstring(gomock.Any(), substring).Return(nil, expectedErr)
		s := rpc.New(mockRepo, env)

		data, err := s.SearchPeers(ctx, &community_proto.SearchPeersIn{Substring: substring})
		assert.Nil(t, data)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "select err")
	})
}
