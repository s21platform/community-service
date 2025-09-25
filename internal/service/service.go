package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/pkg/community"
)

type Service struct {
	community.UnimplementedCommunityServiceServer
	dbR   DbRepo
	env   string
	rR    RedisRepo
	notCl NotificationS
}

func New(dbR DbRepo, env string, rR RedisRepo, notCl NotificationS, cfg *config.Config) *Service {
	return &Service{
		dbR:   dbR,
		env:   env,
		rR:    rR,
		notCl: notCl,
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
