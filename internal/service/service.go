package service

import (
	"context"
	"database/sql"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"strconv"

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

func (s *Service) IsUserStaff(ctx context.Context, in *community.LoginIn) (*community.IsUserStaffOut, error) {
	_, err := s.dbR.GetStaffId(ctx, in.Login)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger_lib.Error(logger_lib.WithError(ctx, err), "cannot check is user staff")
			return nil, status.Errorf(codes.Internal, "cannot check is user staff, err: %v", err)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return &community.IsUserStaffOut{IsStaff: false}, nil
		}
	}

	return &community.IsUserStaffOut{IsStaff: true}, nil
}

func (s *Service) GetPeerSchoolData(ctx context.Context, in *community.GetSchoolDataIn) (*community.GetSchoolDataOut, error) {
	schoolData, err := s.dbR.GetPeerSchoolData(ctx, in.NickName)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "cannot get school data")
		return nil, status.Errorf(codes.Internal, "cannot get peer school data, err: %s", err)
	}
	return &community.GetSchoolDataOut{ClassName: schoolData.ClassName, ParallelName: schoolData.ParallelName}, nil
}

func (s *Service) SearchPeers(ctx context.Context, in *community.SearchPeersIn) (*community.SearchPeersOut, error) {
	log.Println("Input SearchPeers: ", in)
	res, err := s.dbR.SearchPeersBySubstring(ctx, in.Substring)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "search peer error: %s", err)
	}
	return &community.SearchPeersOut{SearchPeers: res}, nil
}

func (s *Service) RunLoginsWorkerManually(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	s.rR.Delete(ctx, config.KeyLoginsLastUpdated)
	return &emptypb.Empty{}, nil
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
		return &community.ValidateCodeOut{Message: "Код не найден"}, nil
	}
	codeInt, err := strconv.Atoi(code)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to convert code to int")
		return &community.ValidateCodeOut{Message: ""}, status.Errorf(codes.Internal, "failed to convert code: %v", err)
	}
	if int64(codeInt) != in.Code {
		return &community.ValidateCodeOut{Message: "Не совпадает код"}, nil
	}
	id, err := s.dbR.GetIdFromParticipant(ctx, uuid)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get user id")
		return &community.ValidateCodeOut{Message: ""}, status.Errorf(codes.NotFound, "failed to get user id, err: %v", err)
	}
	err = s.dbR.InsertLinkEdu(ctx, id, uuid)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to insert link")
		return &community.ValidateCodeOut{Message: ""}, status.Errorf(codes.NotFound, "failed to insert link edu, err: %v", err)
	}
	return nil, nil
}
