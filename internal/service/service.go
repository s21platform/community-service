package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

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
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("IsUserStaff")

	_, err := s.dbR.GetStaffId(ctx, in.Login)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.Error(fmt.Sprintf("cannot check is user staff, err: %v", err))
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
		log.Printf("cannot get peer school data, err: %s\n", err)
		return nil, status.Errorf(codes.Internal, "cannot get peer school data, err: %s", err)
	}
	return &community.GetSchoolDataOut{ClassName: schoolData.ClassName, ParallelName: schoolData.ParallelName}, nil
}

func (s *Service) IsPeerExist(ctx context.Context, in *community.EmailIn) (*community.EmailOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("IsPeerExist")

	peerStatus, err := s.dbR.GetPeerStatus(ctx, in.Email)
	if err != nil {
		logger.Error(fmt.Sprintf("cannot get peer status, err: %v", err))
		return nil, status.Errorf(codes.Internal, "cannot get peer error: %v", err)
	}

	if peerStatus != "ACTIVE" {
		logger.Info(fmt.Sprintf("peer=%s has status: %s", in.Email, peerStatus))
		return &community.EmailOut{IsExist: false}, nil
	}

	if s.env == "stage" {
		_, err := s.dbR.GetStaffId(ctx, in.Email)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				logger.Error(fmt.Sprintf("cannot check is user staff, err: %v", err))
				return nil, status.Errorf(codes.Internal, "cannot check is user staff, err: %v", err)
			}

			if errors.Is(err, sql.ErrNoRows) {
				logger.Info(fmt.Sprintf("user %s is not allowed to the stage enviroment", in.Email))
				return nil, status.Errorf(codes.PermissionDenied, "user %s is not allowed to the stage environment", in.Email)
			}
		}
	}

	return &community.EmailOut{IsExist: true}, nil
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

func (s *Service) SendEduLinkingCode(ctx context.Context, in *community.SendEduLinkingCodeIn) (*emptypb.Empty, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("SendEduLinkingCode")

	peerStatus, err := s.dbR.GetPeerStatus(ctx, in.Login)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get peer status, err: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to get peer status, err: %v", err)
	}

	if peerStatus != "ACTIVE" {
		logger.Info(fmt.Sprintf("peer=%s has status: %s", in.Login, peerStatus))
		return &emptypb.Empty{}, nil
	}

	code := strconv.Itoa(rand.Intn(89999) + 10000)

	err = s.rR.Set(ctx, config.Key("code_"+in.Login), code, time.Minute*10)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to set code to redis, err: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to set code to redis, err: %v", err)
	}

	email := in.Login + "@student.21-school.ru"

	err = s.notCl.SendEduCode(ctx, email, code)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to send verification code, err: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to send verification code, err: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) GetStudentData(ctx context.Context, in *community.GetStudentDataIn) (*community.GetStudentDataOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetStudentData")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to not found UUID in context")
		return nil, status.Error(codes.Internal, "uuid not found in context")
	}

	selfID, err := s.dbR.GetIdPeer(ctx, uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get user id, err: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get user id, err: %v", err)
	}
	peerID, err := s.dbR.GetIdPeer(ctx, in.UserUUID)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get user id, err: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get user id, err: %v", err)
	}
	if selfID == 0 || peerID == 0 {
		logger.Error(fmt.Sprintf("one of the peers is not in the list, err: %v", err))
		return nil, status.Errorf(codes.NotFound, "one of the peers is not in the list, err: %v", err)
	}

	data, err := s.dbR.GetPeerData(ctx, peerID)
	if err != nil {
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
		tmp := &community.Skill{
			Name:   j.Name,
			Points: j.Points,
		}
		out.Skills[i] = tmp
	}

	out.Badges = make([]*community.Badge, len(data.Badges))
	for i, j := range data.Badges {
		tmp := &community.Badge{
			Name:            j.Name,
			ReceiptDateTime: j.ReceiptDateTime,
			IconUrl:         j.IconURL,
		}
		out.Badges[i] = tmp
	}

	return out, nil
}
