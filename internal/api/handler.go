package api

import (
	"encoding/json"
	"github.com/s21platform/community-service/internal/model"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/community-service/internal/config"
	apigen "github.com/s21platform/community-service/internal/generated"
)

type Handler struct {
	dbR DbRepo
	rR  RedisRepo
	nC  NotificationClient
}

func New(dbR DbRepo, rR RedisRepo, nC NotificationClient) *Handler {
	return &Handler{
		dbR: dbR,
		rR:  rR,
		nC:  nC,
	}
}

func (h *Handler) SendLinkingCode(w http.ResponseWriter, r *http.Request, params apigen.SendLinkingCodeParams) {
	ctx := logger_lib.WithUserUuid(r.Context(), params.XUserUuid)

	var body apigen.SendLinkingCodeData
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to decode request body")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	ctx = logger_lib.WithField(ctx, "login", body.Login)

	peerStatus, err := h.dbR.GetPeerStatus(ctx, body.Login)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get peer status")
		resolveError(&w, http.StatusInternalServerError)
		return
	}
	ctx = logger_lib.WithField(ctx, "perr_status", peerStatus)

	if peerStatus != model.ParticipantStatusActive {
		logger_lib.Info(ctx, "peer have inactive status")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	code := strconv.Itoa(rand.Intn(89999) + 10000)

	err = h.rR.Set(ctx, config.Key("code_"+body.Login), code, time.Minute*10)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to set code")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	email := body.Login + "@student.21-school.ru"

	err = h.nC.SendEduCode(ctx, email, code)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to send verification code")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func resolveError(w *http.ResponseWriter, status int) {
	var message string
	switch status {
	case http.StatusBadRequest:
		message = "Произошла ошибка, попробуйте перезагрузить страницу"
	case http.StatusUnauthorized:
		message = "Вы не авторизованы для этого действия"
	case http.StatusNotFound:
		message = "Страница не найдена"
	case http.StatusInternalServerError:
		message = "У нас что-то сломалось, но мы уже чиним!"
	default:
		message = "У нас что-то сломалось, но мы уже чиним!"
	}

	body, _ := json.Marshal(apigen.Forbidden{
		Message: message,
	})
	(*w).WriteHeader(status)
	_, _ = (*w).Write(body)
}
