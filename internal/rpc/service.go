package rpc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	logger_lib "github.com/s21platform/logger-lib"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/pkg/community"
)

type Service struct {
	community.UnimplementedCommunityServiceServer
	dbR DbRepo
	env string
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

func New(dbR DbRepo, env string) *Service {
	return &Service{
		dbR: dbR,
		env: env,
	}
}
