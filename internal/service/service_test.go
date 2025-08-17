package service

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/s21platform/community-service/pkg/community"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/model"
)

var env = "prod"

func TestService_IsUserStaff(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := NewMockDbRepo(controller)

	mockLogger := logger_lib.NewMockLoggerInterface(controller)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	t.Run("is_user_staff_ok", func(t *testing.T) {
		login := "staff_login"
		var id int64 = 1

		mockLogger.EXPECT().AddFuncName("IsUserStaff")
		mockRepo.EXPECT().GetStaffId(gomock.Any(), login).Return(id, nil)

		s := New(mockRepo, env, nil, nil, nil)
		if s == nil {
			t.Errorf("New returned a nil service object")
		}
		if s.dbR == nil {
			t.Errorf("New did not initialize the dbR field")
		}

		data, err := s.IsUserStaff(ctx, &community.LoginIn{Login: login})
		assert.NoError(t, err)
		assert.True(t, data.IsStaff)
	})

	t.Run("is_user_staff_false", func(t *testing.T) {
		login := "not_staff_login"
		var id int64 = 0

		mockLogger.EXPECT().AddFuncName("IsUserStaff")
		mockRepo.EXPECT().GetStaffId(gomock.Any(), login).Return(id, sql.ErrNoRows)

		s := New(mockRepo, env, nil, nil, nil)

		data, err := s.IsUserStaff(ctx, &community.LoginIn{Login: login})
		assert.NoError(t, err)
		assert.False(t, data.IsStaff)
	})

	t.Run("is_user_staff_err", func(t *testing.T) {
		login := "not_staff_login"
		var id int64 = 0
		expectedErr := errors.New("select err")

		mockLogger.EXPECT().AddFuncName("IsUserStaff")
		mockLogger.EXPECT().Error(fmt.Sprintf("cannot check is user staff, err: %v", expectedErr))
		mockRepo.EXPECT().GetStaffId(gomock.Any(), login).Return(id, expectedErr)

		s := New(mockRepo, env, nil, nil, nil)

		data, err := s.IsUserStaff(ctx, &community.LoginIn{Login: login})
		assert.Nil(t, data)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "cannot check is user staff")
	})
}

func TestServer_GetPeerSchoolData(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepo := NewMockDbRepo(controller)

	t.Run("get_peer_school_data_ok", func(t *testing.T) {
		expectedData := model.PeerSchoolData{ClassName: "test-class", ParallelName: "test-parallel"}
		nickName := "aboba"
		mockRepo.EXPECT().GetPeerSchoolData(gomock.Any(), nickName).Return(expectedData, nil)

		s := New(mockRepo, env, nil, nil, nil)
		data, err := s.GetPeerSchoolData(ctx, &community.GetSchoolDataIn{NickName: nickName})
		assert.NoError(t, err)
		assert.Equal(t, data, &community.GetSchoolDataOut{ClassName: expectedData.ClassName, ParallelName: expectedData.ParallelName})
	})

	t.Run("get_peer_school_data_err", func(t *testing.T) {
		nickName := "aboba"
		expectedErr := errors.New("select err")
		mockRepo.EXPECT().GetPeerSchoolData(gomock.Any(), nickName).Return(model.PeerSchoolData{}, expectedErr)

		s := New(mockRepo, env, nil, nil, nil)

		data, err := s.GetPeerSchoolData(ctx, &community.GetSchoolDataIn{NickName: nickName})
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
	mockRepo := NewMockDbRepo(controller)

	mockLogger := logger_lib.NewMockLoggerInterface(controller)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	t.Run("get_peer_status_ok", func(t *testing.T) {
		expectedStatus := "ACTIVE"
		email := "aboba@student.21-school.ru"

		mockRepo.EXPECT().GetPeerStatus(gomock.Any(), email).Return(expectedStatus, nil)
		mockLogger.EXPECT().AddFuncName("IsPeerExist")

		s := New(mockRepo, env, nil, nil, nil)
		isExist, err := s.IsPeerExist(ctx, &community.EmailIn{Email: email})
		assert.NoError(t, err)
		assert.True(t, isExist.IsExist)
	})

	t.Run("get_peer_status_stage_ok", func(t *testing.T) {
		expectedStatus := "ACTIVE"
		email := "aboba@student.21-school.ru"
		var id int64 = 5

		mockRepo.EXPECT().GetPeerStatus(gomock.Any(), email).Return(expectedStatus, nil)
		mockRepo.EXPECT().GetStaffId(gomock.Any(), email).Return(id, nil)
		mockLogger.EXPECT().AddFuncName("IsPeerExist")

		s := New(mockRepo, "stage", nil, nil, nil)
		isExist, err := s.IsPeerExist(ctx, &community.EmailIn{Email: email})
		assert.NoError(t, err)
		assert.True(t, isExist.IsExist)
	})

	// user has no permission to stage cfg
	t.Run("get_peer_status_no_permission_to_stage", func(t *testing.T) {
		expectedStatus := "ACTIVE"
		email := "aboba@student.21-school.ru"
		var id int64 = 0

		mockRepo.EXPECT().GetPeerStatus(gomock.Any(), email).Return(expectedStatus, nil)
		mockRepo.EXPECT().GetStaffId(gomock.Any(), email).Return(id, sql.ErrNoRows)
		mockLogger.EXPECT().AddFuncName("IsPeerExist")
		mockLogger.EXPECT().Info(fmt.Sprintf("user %s is not allowed to the stage enviroment", email))

		s := New(mockRepo, "stage", nil, nil, nil)
		isExist, err := s.IsPeerExist(ctx, &community.EmailIn{Email: email})
		assert.Nil(t, isExist)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.PermissionDenied, st.Code())
		assert.Contains(t, st.Message(), "user aboba@student.21-school.ru is not allowed to the stage environment")
	})

	t.Run("get_peer_status_stage_err", func(t *testing.T) {
		expectedStatus := "ACTIVE"
		email := "aboba@student.21-school.ru"
		var id int64 = 0
		expectedErr := errors.New("select err")

		mockRepo.EXPECT().GetPeerStatus(gomock.Any(), email).Return(expectedStatus, nil)
		mockRepo.EXPECT().GetStaffId(gomock.Any(), email).Return(id, expectedErr)
		mockLogger.EXPECT().AddFuncName("IsPeerExist")
		mockLogger.EXPECT().Error(fmt.Sprintf("cannot check is user staff, err: %v", expectedErr))

		s := New(mockRepo, "stage", nil, nil, nil)
		isExist, err := s.IsPeerExist(ctx, &community.EmailIn{Email: email})
		assert.Nil(t, isExist)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "select err")
	})

	t.Run("get_peer_status_not_found", func(t *testing.T) {
		expectedStatus := "NOT ACTIVE"
		email := "aboba@student.21-school.ru"

		mockRepo.EXPECT().GetPeerStatus(gomock.Any(), email).Return(expectedStatus, nil)
		mockLogger.EXPECT().AddFuncName("IsPeerExist")
		mockLogger.EXPECT().Info(fmt.Sprintf("peer=%s has status: %s", email, expectedStatus))

		s := New(mockRepo, env, nil, nil, nil)
		isExist, err := s.IsPeerExist(ctx, &community.EmailIn{Email: email})
		assert.NoError(t, err)
		assert.False(t, isExist.IsExist)
	})

	t.Run("get_peer_status_err", func(t *testing.T) {
		email := "aboba@student.21-school.ru"
		expectedErr := errors.New("select err")

		mockRepo.EXPECT().GetPeerStatus(gomock.Any(), email).Return("", expectedErr)
		mockLogger.EXPECT().AddFuncName("IsPeerExist")
		mockLogger.EXPECT().Error(fmt.Sprintf("cannot get peer status, err: %v", expectedErr))

		s := New(mockRepo, env, nil, nil, nil)
		isExist, err := s.IsPeerExist(ctx, &community.EmailIn{Email: email})

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
	mockRepo := NewMockDbRepo(controller)

	t.Run("search_peers_ok", func(t *testing.T) {
		expectedData := []*community.SearchPeer{
			{Login: "aboba"},
			{Login: "abobaoba"},
			{Login: "aboo"},
		}
		substring := "ab"
		mockRepo.EXPECT().SearchPeersBySubstring(gomock.Any(), substring).Return(expectedData, nil)

		s := New(mockRepo, env, nil, nil, nil)
		data, err := s.SearchPeers(ctx, &community.SearchPeersIn{Substring: substring})
		assert.NoError(t, err)
		assert.Equal(t, data, &community.SearchPeersOut{SearchPeers: expectedData})
	})

	t.Run("search_peers_err", func(t *testing.T) {
		expectedErr := errors.New("select err")
		substring := "ab"
		mockRepo.EXPECT().SearchPeersBySubstring(gomock.Any(), substring).Return(nil, expectedErr)
		s := New(mockRepo, env, nil, nil, nil)

		data, err := s.SearchPeers(ctx, &community.SearchPeersIn{Substring: substring})
		assert.Nil(t, data)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "select err")
	})
}

func TestServer_SendEduLinkingCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepo := NewMockDbRepo(controller)
	mockRedisRepo := NewMockRedisRepo(controller)
	mockLogger := logger_lib.NewMockLoggerInterface(controller)
	mockNotCl := NewMockNotificationS(controller)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	t.Run("send_code_email_ok", func(t *testing.T) {
		login := "aboba"
		email := login + "@student.21-school.ru"
		mockRepo.EXPECT().GetPeerStatus(gomock.Any(), email).Return("ACTIVE", nil)
		mockRedisRepo.EXPECT().Set(gomock.Any(), config.Key("code_"+login), gomock.Any(), gomock.Any()).Return(nil)
		mockLogger.EXPECT().AddFuncName("SendEduLinkingCode")
		mockNotCl.EXPECT().SendVerificationCode(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		s := New(mockRepo, env, mockRedisRepo, mockNotCl, nil)
		_, err := s.SendEduLinkingCode(ctx, &community.SendEduLinkingCodeIn{Login: login})
		assert.NoError(t, err)
	})

	t.Run("send_code_email_err", func(t *testing.T) {
		login := "aboba"
		email := login + "@student.21-school.ru"
		expectedErr := errors.New("set err")
		mockRepo.EXPECT().GetPeerStatus(gomock.Any(), email).Return("ACTIVE", nil)
		mockRedisRepo.EXPECT().Set(gomock.Any(), config.Key("code_"+login), gomock.Any(), gomock.Any()).Return(expectedErr)
		mockLogger.EXPECT().AddFuncName("SendEduLinkingCode")
		mockLogger.EXPECT().Error(fmt.Sprintf("failed to set code to redis, err: %v", expectedErr))

		s := New(mockRepo, env, mockRedisRepo, mockNotCl, nil)
		_, err := s.SendEduLinkingCode(ctx, &community.SendEduLinkingCodeIn{Login: login})
		assert.Error(t, err)
	})
}
