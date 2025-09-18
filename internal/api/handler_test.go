package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	apigen "github.com/s21platform/community-service/internal/generated"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/stretchr/testify/assert"
)

func TestHandler_SendLinkingCode(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockDbRepo := NewMockDbRepo(controller)
	mockRedisRepo := NewMockRedisRepo(controller)
	mockNotificationClient := NewMockNotificationClient(controller)

	handler := New(mockDbRepo, mockRedisRepo, mockNotificationClient)

	t.Run("success_case", func(t *testing.T) {
		userUUID := "test-uuid"
		login := "test_login"
		requestBody := apigen.SendLinkingCodeData{
			Login: login,
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/link", bytes.NewReader(bodyBytes))
		req = req.WithContext(logger_lib.WithUserUuid(context.Background(), userUUID))
		w := httptest.NewRecorder()

		mockDbRepo.EXPECT().
			GetPeerStatus(gomock.Any(), login).
			Return("ACTIVE", nil)

		mockRedisRepo.EXPECT().
			Set(gomock.Any(), gomock.Any(), gomock.Any(), time.Minute*10).
			Return(nil)

		mockNotificationClient.EXPECT().
			SendEduCode(gomock.Any(), login+"@student.21-school.ru", gomock.Any()).
			Return(nil)

		handler.SendLinkingCode(w, req, apigen.SendLinkingCodeParams{XUserUuid: userUUID})

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid_request_body", func(t *testing.T) {
		userUUID := "test-uuid"
		invalidBody := []byte("invalid json")

		req := httptest.NewRequest(http.MethodPost, "/link", bytes.NewReader(invalidBody))
		req = req.WithContext(logger_lib.WithUserUuid(context.Background(), userUUID))
		w := httptest.NewRecorder()

		handler.SendLinkingCode(w, req, apigen.SendLinkingCodeParams{XUserUuid: userUUID})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response apigen.Forbidden
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "У нас что-то сломалось, но мы уже чиним!", response.Message)
	})

	t.Run("get_peer_status_error", func(t *testing.T) {
		userUUID := "test-uuid"
		login := "test_login"
		requestBody := apigen.SendLinkingCodeData{
			Login: login,
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/link", bytes.NewReader(bodyBytes))
		req = req.WithContext(logger_lib.WithUserUuid(context.Background(), userUUID))
		w := httptest.NewRecorder()

		mockDbRepo.EXPECT().
			GetPeerStatus(gomock.Any(), login).
			Return("", errors.New("db error"))

		handler.SendLinkingCode(w, req, apigen.SendLinkingCodeParams{XUserUuid: userUUID})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("inactive_peer_status", func(t *testing.T) {
		userUUID := "test-uuid"
		login := "test_login"
		requestBody := apigen.SendLinkingCodeData{
			Login: login,
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/link", bytes.NewReader(bodyBytes))
		req = req.WithContext(logger_lib.WithUserUuid(context.Background(), userUUID))
		w := httptest.NewRecorder()

		mockDbRepo.EXPECT().
			GetPeerStatus(gomock.Any(), login).
			Return("INACTIVE", nil)

		handler.SendLinkingCode(w, req, apigen.SendLinkingCodeParams{XUserUuid: userUUID})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("redis_set_error", func(t *testing.T) {
		userUUID := "test-uuid"
		login := "test_login"
		requestBody := apigen.SendLinkingCodeData{
			Login: login,
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/link", bytes.NewReader(bodyBytes))
		req = req.WithContext(logger_lib.WithUserUuid(context.Background(), userUUID))
		w := httptest.NewRecorder()

		mockDbRepo.EXPECT().
			GetPeerStatus(gomock.Any(), login).
			Return("ACTIVE", nil)

		mockRedisRepo.EXPECT().
			Set(gomock.Any(), gomock.Any(), gomock.Any(), time.Minute*10).
			Return(errors.New("redis error"))

		handler.SendLinkingCode(w, req, apigen.SendLinkingCodeParams{XUserUuid: userUUID})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("notification_send_error", func(t *testing.T) {
		userUUID := "test-uuid"
		login := "test_login"
		requestBody := apigen.SendLinkingCodeData{
			Login: login,
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/link", bytes.NewReader(bodyBytes))
		req = req.WithContext(logger_lib.WithUserUuid(context.Background(), userUUID))
		w := httptest.NewRecorder()

		mockDbRepo.EXPECT().
			GetPeerStatus(gomock.Any(), login).
			Return("ACTIVE", nil)

		mockRedisRepo.EXPECT().
			Set(gomock.Any(), gomock.Any(), gomock.Any(), time.Minute*10).
			Return(nil)

		mockNotificationClient.EXPECT().
			SendEduCode(gomock.Any(), login+"@student.21-school.ru", gomock.Any()).
			Return(errors.New("notification error"))

		handler.SendLinkingCode(w, req, apigen.SendLinkingCodeParams{XUserUuid: userUUID})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestResolveError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		statusCode     int
		expectedMsg    string
		expectedStatus int
	}{
		{
			name:           "bad_request",
			statusCode:     http.StatusBadRequest,
			expectedMsg:    "Произошла ошибка, попробуйте перезагрузить страницу",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "unauthorized",
			statusCode:     http.StatusUnauthorized,
			expectedMsg:    "Вы не авторизованы для этого действия",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "not_found",
			statusCode:     http.StatusNotFound,
			expectedMsg:    "Страница не найдена",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "internal_server_error",
			statusCode:     http.StatusInternalServerError,
			expectedMsg:    "У нас что-то сломалось, но мы уже чиним!",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "unknown_error",
			statusCode:     http.StatusTeapot,
			expectedMsg:    "У нас что-то сломалось, но мы уже чиним!",
			expectedStatus: http.StatusTeapot,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			rw := http.ResponseWriter(w)
			resolveError(&rw, tc.statusCode)

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response apigen.Forbidden
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedMsg, response.Message)
		})
	}
}
