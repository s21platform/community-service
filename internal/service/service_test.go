package service

import (
	"errors"
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
	"github.com/s21platform/community-service/internal/pkg/tx"
)

var env = "prod"

type noopTxRepo struct{}

func (noopTxRepo) WithTx(ctx context.Context, cb func(ctx context.Context) error) error {
	return cb(ctx)
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

		s := New(mockRepo, env, nil, nil, nil, nil)
		data, err := s.GetPeerSchoolData(ctx, &community.GetSchoolDataIn{NickName: nickName})
		assert.NoError(t, err)
		assert.Equal(t, data, &community.GetSchoolDataOut{ClassName: expectedData.ClassName, ParallelName: expectedData.ParallelName})
	})

	t.Run("get_peer_school_data_err", func(t *testing.T) {
		nickName := "aboba"
		expectedErr := errors.New("select err")
		mockRepo.EXPECT().GetPeerSchoolData(gomock.Any(), nickName).Return(model.PeerSchoolData{}, expectedErr)

		s := New(mockRepo, env, nil, nil, nil, nil)

		data, err := s.GetPeerSchoolData(ctx, &community.GetSchoolDataIn{NickName: nickName})
		assert.Nil(t, data)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "select err")
	})
}

func TestService_GetStudentData(t *testing.T) {
	t.Parallel()
	ctx := context.WithValue(context.Background(), config.KeyUUID, "uuid-1")
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepo := NewMockDbRepo(controller)
	mockRedisRepo := NewMockRedisRepo(controller)
	mockNotCl := NewMockNotificationS(controller)
	mockLogger := logger_lib.NewMockLoggerInterface(controller)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	t.Run("success_case", func(t *testing.T) {
		inputUUID := "user-2"
		ctxUUID := "uuid-1"
		request := &community.GetStudentDataIn{UserUUID: inputUUID}
		var idFirst int64 = 1
		var idSecond int64 = 2

		mockRepo.EXPECT().
			GetIdPeer(ctx, ctxUUID).
			Return(idFirst, nil).
			Times(1)
		mockRepo.EXPECT().
			GetIdPeer(ctx, inputUUID).
			Return(idSecond, nil).
			Times(1)
		mockData := &model.ParticipantData{
			Login:          "login1",
			CampusId:       1,
			ClassName:      "10A",
			ParallelName:   "10",
			TribeID:        5,
			Status:         "active",
			CreatedAt:      "2025-08-16",
			ExpValue:       1000,
			Level:          5,
			ExpToNextLevel: 200,
			Crp:            10,
			Prp:            20,
			Coins:          50,
			Skills:         []model.Skill{{Name: "skill1", Points: 100}},
			Badges:         []model.Badge{{Name: "badge1", ReceiptDateTime: "2025-08-16", IconURL: "url"}},
		}
		mockRepo.EXPECT().
			GetPeerData(ctx, idSecond).
			Return(mockData, nil).
			Times(1)

		s := New(mockRepo, env, mockRedisRepo, mockNotCl, nil, nil)
		_, err := s.GetStudentData(ctx, request)
		assert.NoError(t, err)
	})

	t.Run("get_first_id_error", func(t *testing.T) {
		inputUUID := "user-2"
		ctxUUID := "uuid-1"
		request := &community.GetStudentDataIn{UserUUID: inputUUID}
		expectedErr := errors.New("get id error")

		mockRepo.EXPECT().
			GetIdPeer(ctx, ctxUUID).
			Return(int64(0), expectedErr).
			Times(1)

		s := New(mockRepo, env, mockRedisRepo, mockNotCl, nil, nil)
		data, err := s.GetStudentData(ctx, request)
		assert.Nil(t, data)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("get_second_id_error", func(t *testing.T) {
		inputUUID := "user-2"
		ctxUUID := "uuid-1"
		request := &community.GetStudentDataIn{UserUUID: inputUUID}
		expectedErr := errors.New("get id error")

		mockRepo.EXPECT().
			GetIdPeer(ctx, ctxUUID).
			Return(int64(1), nil).
			Times(1)
		mockRepo.EXPECT().
			GetIdPeer(ctx, inputUUID).
			Return(int64(0), expectedErr).
			Times(1)

		s := New(mockRepo, env, mockRedisRepo, mockNotCl, nil, nil)
		data, err := s.GetStudentData(ctx, request)
		assert.Nil(t, data)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("get_peer_data_error", func(t *testing.T) {
		inputUUID := "user-2"
		ctxUUID := "uuid-1"
		request := &community.GetStudentDataIn{UserUUID: inputUUID}
		expectedErr := errors.New("get peer data error")

		mockRepo.EXPECT().
			GetIdPeer(ctx, ctxUUID).
			Return(int64(1), nil).
			Times(1)
		mockRepo.EXPECT().
			GetIdPeer(ctx, inputUUID).
			Return(int64(2), nil).
			Times(1)
		mockRepo.EXPECT().
			GetPeerData(ctx, int64(2)).
			Return(nil, expectedErr).
			Times(1)

		s := New(mockRepo, env, mockRedisRepo, mockNotCl, nil, nil)
		data, err := s.GetStudentData(ctx, request)
		assert.Nil(t, data)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
	})
}

func TestService_ValidateCode(t *testing.T) {
	t.Parallel()
	baseCtx := context.WithValue(context.Background(), config.KeyUUID, "uuid-1")
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepo := NewMockDbRepo(controller)
	mockRedisRepo := NewMockRedisRepo(controller)
	mockNotCl := NewMockNotificationS(controller)
	mockLogger := logger_lib.NewMockLoggerInterface(controller)
	baseCtx = context.WithValue(baseCtx, config.KeyLogger, mockLogger)
	baseCtx = context.WithValue(baseCtx, tx.KeyTx, tx.Tx{DbRepo: noopTxRepo{}})

	t.Run("success_case", func(t *testing.T) {
		ctx := baseCtx
		key := "15"
		ctxUUID := "uuid-1"
		var id int64 = 15
		login := "test1"
		request := &community.ValidateCodeIn{Login: login, Code: 15}

		mockRedisRepo.EXPECT().
			GetByKey(ctx, gomock.Any()).
			Return(key, nil).
			Times(1)
		mockRepo.EXPECT().
			GetIdFromParticipant(ctx, ctxUUID).
			Return(id, nil).
			Times(1)
		mockRepo.EXPECT().
			InsertLinkEdu(ctx, id, ctxUUID).
			Return(nil).
			Times(1)
		mockRepo.EXPECT().
			GetLogin(ctx, id).
			Return(login, nil).
			Times(1)

		mockUlE := NewMockUserLinkingEdu(controller)
		mockUlE.EXPECT().
			ProduceMessage(ctx, gomock.Any(), ctxUUID).
			Return(nil).
			Times(1)

		s := New(mockRepo, env, mockRedisRepo, mockNotCl, mockUlE, nil)
		_, err := s.ValidateCode(ctx, request)
		assert.NoError(t, err)
	})

	t.Run("GetByKey_err", func(t *testing.T) {
		ctx := baseCtx
		expectedErr := errors.New("get err")
		request := &community.ValidateCodeIn{Login: "test1", Code: 15}

		mockRedisRepo.EXPECT().
			GetByKey(ctx, gomock.Any()).
			Return("", expectedErr).
			Times(1)

		s := New(mockRepo, env, mockRedisRepo, mockNotCl, nil, nil)
		_, err := s.ValidateCode(ctx, request)
		assert.Error(t, err)
	})

	t.Run("code_err", func(t *testing.T) {
		ctx := baseCtx
		request := &community.ValidateCodeIn{Login: "test1", Code: 15}

		mockRedisRepo.EXPECT().
			GetByKey(ctx, gomock.Any()).
			Return("", nil).
			Times(1)

		s := New(mockRepo, env, mockRedisRepo, mockNotCl, nil, nil)
		_, err := s.ValidateCode(ctx, request)
		assert.NoError(t, err)
	})

	t.Run("atoi_err", func(t *testing.T) {
		ctx := baseCtx
		key := "test"
		request := &community.ValidateCodeIn{Login: "test1", Code: 15}

		mockRedisRepo.EXPECT().
			GetByKey(ctx, gomock.Any()).
			Return(key, nil).
			Times(1)

		s := New(mockRepo, env, mockRedisRepo, mockNotCl, nil, nil)
		_, err := s.ValidateCode(ctx, request)
		assert.Error(t, err)
	})

	t.Run("get_id_from_participant_err", func(t *testing.T) {
		ctx := baseCtx
		key := "15"
		ctxUUID := "uuid-1"
		request := &community.ValidateCodeIn{Login: "test1", Code: 15}
		expectedErr := errors.New("get id error")

		mockRedisRepo.EXPECT().
			GetByKey(ctx, gomock.Any()).
			Return(key, nil).
			Times(1)
		mockRepo.EXPECT().
			GetIdFromParticipant(ctx, ctxUUID).
			Return(int64(0), expectedErr).
			Times(1)

		s := New(mockRepo, env, mockRedisRepo, mockNotCl, nil, nil)
		_, err := s.ValidateCode(ctx, request)
		assert.Error(t, err)
	})

	t.Run("insert_link_edu_err", func(t *testing.T) {
		ctx := baseCtx
		key := "15"
		ctxUUID := "uuid-1"
		var id int64 = 15
		request := &community.ValidateCodeIn{Login: "test1", Code: 15}
		expectedErr := errors.New("insert error")

		mockRedisRepo.EXPECT().
			GetByKey(ctx, gomock.Any()).
			Return(key, nil).
			Times(1)
		mockRepo.EXPECT().
			GetIdFromParticipant(ctx, ctxUUID).
			Return(id, nil).
			Times(1)
		mockRepo.EXPECT().
			InsertLinkEdu(ctx, id, ctxUUID).
			Return(expectedErr).
			Times(1)

		s := New(mockRepo, env, mockRedisRepo, mockNotCl, nil, nil)
		_, err := s.ValidateCode(ctx, request)
		assert.Error(t, err)
	})

	t.Run("get_login_err", func(t *testing.T) {
		ctx := baseCtx
		key := "15"
		ctxUUID := "uuid-1"
		var id int64 = 15
		request := &community.ValidateCodeIn{Login: "test1", Code: 15}
		expectedErr := errors.New("failed to get login")

		mockRedisRepo.EXPECT().
			GetByKey(ctx, gomock.Any()).
			Return(key, nil).
			Times(1)
		mockRepo.EXPECT().
			GetIdFromParticipant(ctx, ctxUUID).
			Return(id, nil).
			Times(1)
		mockRepo.EXPECT().
			InsertLinkEdu(ctx, id, ctxUUID).
			Return(nil).
			Times(1)
		mockRepo.EXPECT().
			GetLogin(ctx, id).
			Return("", expectedErr).
			Times(1)

		s := New(mockRepo, env, mockRedisRepo, mockNotCl, nil, nil)
		_, err := s.ValidateCode(ctx, request)
		assert.Error(t, err)
	})

	t.Run("ctx_err", func(t *testing.T) {
		ctx := context.Background()
		request := &community.ValidateCodeIn{Login: "test1", Code: 15}

		s := New(mockRepo, env, mockRedisRepo, mockNotCl, nil, nil)
		_, err := s.ValidateCode(ctx, request)
		assert.Error(t, err)
	})

}
