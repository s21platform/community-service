package peerdata

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/model"
	"github.com/s21platform/community-service/pkg/community"

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

	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("ParticipantDataWorker")
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("participant uploading worker shutting down")
			return

		case <-ticker.C:
			lastUpdate, err := s.rR.GetByKey(ctx, config.KeyParticipantDataLastUpdated)
			if err != nil {
				logger.Error(fmt.Sprintf("failed to get last update time, err: %v", err))
			}

			if lastUpdate == "" {
				err := s.uploadDataParticipant(ctx)
				if err != nil {
					logger.Error(fmt.Sprintf("failed to upload participants, err: %v", err))
				}

				// по сути мы тут указываем через сколько запустить следующий цикл опроса. 5 часов много, поставил 10 минут передышки
				err = s.rR.Set(ctx, config.KeyParticipantDataLastUpdated, "upd", 10*time.Minute)
				if err != nil {
					logger.Error(fmt.Sprintf("failed to save participant last updated, err: %v", err))
				}
			}
			logger.Info("participant worker done")
		}
	}
}

func (s *Worker) uploadDataParticipant(ctx context.Context) error {
	var offset int64
	mtx := pkg.FromContext(ctx, config.KeyMetrics)
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
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
			exists := true
			participant, err := s.dbR.ParticipantData(ctx, login)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					mtx.Increment("update_participant_data.error_get_participant")
					logger.Error(fmt.Sprintf("failed to check participant existance: %v", err))
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

				logger.Error(fmt.Sprintf("failed to get participant data for login %s, err: %v", login, err))
				continue
			}

			if participantData == nil {
				mtx.Increment("update_participant_data.not_exists")
				continue
			}

			campus, err := s.dbR.GetCampusByUUID(ctx, participantData.CampusUUID)
			if err != nil {
				mtx.Increment("update_participant_data.error_get_campus")
				logger.Error(fmt.Sprintf("failed to get participant campus: %v", err))
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
				logger.Error(fmt.Sprintf("failed to save participant data for login %s, err: %v", login, err))
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
					logger.Error(fmt.Sprintf("failed to produce level change event for %s: %v", login, err))
				}
			}

			if participant.ExpLevel != participantData.ExpValue {
				event := &community.ParticipantChangeEvent{
					Login:    login,
					OldValue: &community.ParticipantChangeEvent_OldValueInt{OldValueInt: int32(participant.ExpLevel)},
					NewValue: &community.ParticipantChangeEvent_NewValueInt{NewValueInt: int32(participantData.ExpValue)},
					At:       timestamppb.Now(),
				}
				if err := s.elcP.ProduceMessage(ctx, event, login); err != nil {
					logger.Error(fmt.Sprintf("failed to produce exp level change event for %s: %v", login, err))
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
					logger.Error(fmt.Sprintf("failed to produce status change event for %s: %v", login, err))
				}
			}

			mtx.Increment("update_participant_data.ok")
		}

		offset += limit
	}
	mtx.Increment("update_participant_data.finish_upload")
	return nil
}
