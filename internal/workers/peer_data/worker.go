package peerdata

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/model"
	"github.com/s21platform/community-service/pkg/community"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	limit = 10000
)

type Worker struct {
	sC   SchoolC
	dbR  DbRepo
	rR   RedisRepo
	lcP  LevelChangeProducer
	elcP ExpLevelChangeProducer
	scP  StatusChangeProducer
}

func New(school SchoolC, dbR DbRepo, rR RedisRepo, lcP LevelChangeProducer, elcP ExpLevelChangeProducer, scP StatusChangeProducer) *Worker {
	return &Worker{
		sC:   school,
		dbR:  dbR,
		rR:   rR,
		lcP:  lcP,
		elcP: elcP,
		scP:  scP,
	}
}

func (s *Worker) RunParticipantWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger_lib.Info(ctx, "participant uploading worker shutting down")
			return

		case <-ticker.C:
			lastUpdate, err := s.rR.GetByKey(ctx, config.KeyParticipantDataLastUpdated)
			if err != nil {
				logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get last update time")
			}
			if lastUpdate == "" {
				err := s.uploadDataParticipant(ctx)
				if err != nil {
					logger_lib.Error(logger_lib.WithError(ctx, err), "failed to upload data participant")
				}

				//по сути мы тут указываем через сколько запустить следующий цикл опроса. 5 часов много, поставил 10 минут передышки
				err = s.rR.Set(ctx, config.KeyParticipantDataLastUpdated, "upd", 10*time.Minute)
				if err != nil {
					logger_lib.Error(logger_lib.WithError(ctx, err), "failed to set last update time")
				}
			}
		}
	}
}

func (s *Worker) uploadDataParticipant(ctx context.Context) error {
	var offset int64
	mtx := pkg.FromContext(ctx, config.KeyMetrics)

	for {
		logins, err := s.dbR.GetParticipantsLogin(ctx, limit, offset)
		if err != nil {
			mtx.Increment("update_participant_data.edu_error")
			return fmt.Errorf("failed to get participant logins, err: %v", err)
		}
		if len(logins) == 0 {
			mtx.Increment("update_participant_data.empty_login_list")
			break
		}

		for _, login := range logins {
			ctx = logger_lib.WithField(ctx, "login", login)
			exists := true
			participant, err := s.dbR.ParticipantData(ctx, login)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					mtx.Increment("update_participant_data.error_get_participant")
					logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get participant")
					continue
				}
				exists = false
			}
			if participant != nil && participant.Status != model.ParticipantStatusActive {
				mtx.Increment("update_participant_data.skip_not_active")
				continue
			}
			time.Sleep(300 * time.Millisecond)
			participantData, err := s.sC.GetParticipantData(ctx, login)
			if err != nil {
				if strings.Contains(err.Error(), "Invalid token") {
					mtx.Increment("update_participant_data.invalid_token")
				} else if strings.Contains(err.Error(), "Too many requests") {
					mtx.Increment("update_participant_data.too_many_requests")
				} else {
					mtx.Increment("update_participant_data.unknown_error")
				}
				logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get participant")
				continue
			}
			if participantData == nil {
				mtx.Increment("update_participant_data.not_exists")
				continue
			}
			campus, err := s.dbR.GetCampusByUUID(ctx, participantData.CampusUUID)
			if err != nil {
				mtx.Increment("update_participant_data.error_get_campus")
				logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get campus")
				continue
			}
			participantData.TribeID = 1
			if !exists {
				err = s.dbR.InsertParticipantData(ctx, participantData, login, campus.Id)
			} else {
				err = s.dbR.UpdateParticipantData(ctx, participantData, login, campus.Id)
			}
			if err != nil {
				mtx.Increment("update_participant_data.not_save")
				logger_lib.Error(logger_lib.WithError(ctx, err), "failed to save participant")
				continue
			}
			if participant.Level != participantData.Level {
				event := &community.ParticipantChangeEvent{
					Login:    login,
					OldValue: &community.ParticipantChangeEvent_OldValueInt{OldValueInt: int32(participant.Level)},
					NewValue: &community.ParticipantChangeEvent_NewValueInt{NewValueInt: int32(participantData.Level)},
					At:       timestamppb.Now(),
				}
				if err := s.lcP.ProduceMessage(ctx, event, login); err != nil {
					logger_lib.Error(logger_lib.WithError(ctx, err), "failed to produce level changed message")
				}
			}

			if participant.ExpValue != participantData.ExpValue {
				event := &community.ParticipantChangeEvent{
					Login:    login,
					OldValue: &community.ParticipantChangeEvent_OldValueInt{OldValueInt: int32(participant.ExpValue)},
					NewValue: &community.ParticipantChangeEvent_NewValueInt{NewValueInt: int32(participantData.ExpValue)},
					At:       timestamppb.Now(),
				}
				if err := s.elcP.ProduceMessage(ctx, event, login); err != nil {
					logger_lib.Error(logger_lib.WithError(ctx, err), "failed to produce exp level changed message")
				}
			}

			if participant.Status != participantData.Status {
				event := &community.ParticipantChangeEvent{
					Login:    login,
					OldValue: &community.ParticipantChangeEvent_OldValueStr{OldValueStr: participant.Status},
					NewValue: &community.ParticipantChangeEvent_NewValueStr{NewValueStr: participantData.Status},
					At:       timestamppb.Now(),
				}
				if err := s.scP.ProduceMessage(ctx, event, login); err != nil {
					logger_lib.Error(logger_lib.WithError(ctx, err), "failed to produce status change message")
				}
			}

			mtx.Increment("update_participant_data.ok")
		}

		offset += limit
	}
	mtx.Increment("update_participant_data.finish_upload")
	return nil
}
