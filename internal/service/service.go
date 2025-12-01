package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/model"
	"github.com/s21platform/community-service/internal/pkg/tx"
	"github.com/s21platform/community-service/pkg/community"
)

type Service struct {
	community.UnimplementedCommunityServiceServer
	dbR   DbRepo
	env   string
	rR    RedisRepo
	notCl NotificationS
	ulE   UserLinkingEdu
}

func New(dbR DbRepo, env string, rR RedisRepo, notCl NotificationS, ulE UserLinkingEdu, cfg *config.Config) *Service {
	return &Service{
		dbR:   dbR,
		env:   env,
		rR:    rR,
		notCl: notCl,
		ulE:   ulE,
	}
}

func (s *Service) GetPeerSchoolData(ctx context.Context, in *community.GetSchoolDataIn) (*community.GetSchoolDataOut, error) {
	schoolData, err := s.dbR.GetPeerSchoolData(ctx, in.NickName)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get school data")
		return nil, status.Errorf(codes.Internal, "failed to get peer school data, err: %s", err)
	}
	return &community.GetSchoolDataOut{ClassName: schoolData.ClassName, ParallelName: schoolData.ParallelName}, nil
}

func (s *Service) GetStudentData(ctx context.Context, in *community.GetStudentDataIn) (*community.GetStudentDataOut, error) {
	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger_lib.Error(ctx, "failed to not found UUID in context")
		return nil, status.Error(codes.Internal, "failed to not found UUID in context")
	}
	// проверка вхождения uuid инициатора в таблицу link_edu
	_, err := s.dbR.GetIdPeer(ctx, uuid)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get peer info")
		return nil, status.Errorf(codes.NotFound, "failed to get user id, err: %v", err)
	}
	// проверка вхождения uuid целевого пользователя в таблицу link_edu
	peerID, err := s.dbR.GetIdPeer(ctx, in.UserUUID)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get peer into linking table")
		return nil, status.Errorf(codes.NotFound, "failed to get user id, err: %v", err)
	}

	data, err := s.dbR.GetPeerData(ctx, peerID)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get peer data")
		return nil, status.Errorf(codes.Internal, "failed to get peer data: %v", err)
	}
	out := &community.GetStudentDataOut{
		Login:          data.Login,
		CampusId:       data.CampusId,
		ClassName:      data.ClassName,
		ParallelName:   data.ParallelName,
		TribeId:        data.TribeID,
		Status:         data.Status,
		CreatedAt:      data.CreatedAt,
		ExpValue:       data.ExpValue,
		Level:          data.Level,
		ExpToNextLevel: data.ExpToNextLevel,
		Crp:            data.Crp,
		Prp:            data.Prp,
		Coins:          data.Coins,
	}
	out.Skills = make([]*community.Skill, len(data.Skills))
	for i, j := range data.Skills {
		out.Skills[i] = &community.Skill{
			Name:   j.Name,
			Points: j.Points,
		}
	}

	out.Badges = make([]*community.Badge, len(data.Badges))
	for i, j := range data.Badges {
		out.Badges[i] = &community.Badge{
			Name:            j.Name,
			ReceiptDateTime: j.ReceiptDateTime,
			IconUrl:         j.IconURL,
		}
	}

	return out, nil
}

func (s *Service) ValidateCode(ctx context.Context, in *community.ValidateCodeIn) (*community.ValidateCodeOut, error) {
	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger_lib.Error(ctx, "failed to not found UUID in context")
		return &community.ValidateCodeOut{Message: ""}, status.Error(codes.Internal, "failed to not found UUID in context")
	}
	code, err := s.rR.GetByKey(ctx, config.Key(in.Login))
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get user code")
		return &community.ValidateCodeOut{Message: ""}, status.Errorf(codes.Internal, "failed to get by key: %v", err)
	}
	if code == "" {
		logger_lib.Error(ctx, "code is not found")
		return &community.ValidateCodeOut{Message: "Код не найден"}, nil
	}
	codeInt, err := strconv.Atoi(code)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to convert code to int")
		return &community.ValidateCodeOut{Message: ""}, status.Errorf(codes.Internal, "failed to convert code: %v", err)
	}
	if int64(codeInt) != in.Code {
		logger_lib.Error(ctx, "failed equal code")
		return &community.ValidateCodeOut{Message: "Не совпадает код"}, nil
	}
	id, err := s.dbR.GetIdFromParticipant(ctx, uuid)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get user id")
		return &community.ValidateCodeOut{Message: ""}, status.Errorf(codes.NotFound, "failed to get user id, err: %v", err)
	}

	err = tx.TxExecute(ctx, func(ctx context.Context) error {
		err = s.dbR.InsertLinkEdu(ctx, id, uuid)
		if err != nil {
			return status.Error(codes.NotFound, "failed to insert link edu, err")
		}

		login, err := s.dbR.GetLogin(ctx, id)
		if err != nil {
			return status.Error(codes.Internal, "failed to get login")
		}

		userLink := model.LinkData{
			UUID:  uuid,
			Login: login,
		}
		rawMessage, err := json.Marshal(userLink)
		if err != nil {
			return fmt.Errorf("failed to marshal user: %v", err)
		}
		err = s.ulE.ProduceMessage(ctx, community.UserCreatedMessage{
			UserUuid:   uuid,
			Login:      login,
			RawMessage: rawMessage,
		}, uuid)
		if err != nil {
			return fmt.Errorf("failed to produce message: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate code: %v", err)
	}
	return &community.ValidateCodeOut{}, nil
}
