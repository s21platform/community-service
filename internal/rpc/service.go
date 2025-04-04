package rpc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	communityproto "github.com/s21platform/community-proto/community-proto"
	logger_lib "github.com/s21platform/logger-lib"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/s21platform/community-service/internal/config"
)

type Service struct {
	communityproto.UnimplementedCommunityServiceServer
	dbR DbRepo
	env string
}

func (s *Service) IsUserStaff(ctx context.Context, in *communityproto.LoginIn) (*communityproto.IsUserStaffOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("IsUserStaff")

	_, err := s.dbR.GetStaffId(ctx, in.Login)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.Error(fmt.Sprintf("cannot check is user staff, err: %v", err))
			return nil, status.Errorf(codes.Internal, "cannot check is user staff, err: %v", err)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return &communityproto.IsUserStaffOut{IsStaff: false}, nil
		}
	}

	return &communityproto.IsUserStaffOut{IsStaff: true}, nil
}

func (s *Service) GetPeerSchoolData(ctx context.Context, in *communityproto.GetSchoolDataIn) (*communityproto.GetSchoolDataOut, error) {
	schoolData, err := s.dbR.GetPeerSchoolData(ctx, in.NickName)
	if err != nil {
		log.Printf("cannot get peer school data, err: %s\n", err)
		return nil, status.Errorf(codes.Internal, "cannot get peer school data, err: %s", err)
	}
	return &communityproto.GetSchoolDataOut{ClassName: schoolData.ClassName, ParallelName: schoolData.ParallelName}, nil
}

func (s *Service) IsPeerExist(ctx context.Context, in *communityproto.EmailIn) (*communityproto.EmailOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("IsPeerExist")

	peerStatus, err := s.dbR.GetPeerStatus(ctx, in.Email)
	if err != nil {
		logger.Error(fmt.Sprintf("cannot get peer status, err: %v", err))
		return nil, status.Errorf(codes.Internal, "cannot get peer error: %v", err)
	}

	if peerStatus != "ACTIVE" {
		logger.Info(fmt.Sprintf("peer=%s has status: %s", in.Email, peerStatus))
		return &communityproto.EmailOut{IsExist: false}, nil
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

	return &communityproto.EmailOut{IsExist: true}, nil
}

func (s *Service) SearchPeers(ctx context.Context, in *communityproto.SearchPeersIn) (*communityproto.SearchPeersOut, error) {
	log.Println("Input SearchPeers: ", in)
	res, err := s.dbR.SearchPeersBySubstring(ctx, in.Substring)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "search peer error: %s", err)
	}
	return &communityproto.SearchPeersOut{SearchPeers: res}, nil
}

func New(dbR DbRepo, env string) *Service {
	return &Service{
		dbR: dbR,
		env: env,
	}
}
