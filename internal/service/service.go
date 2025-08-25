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
	"google.golang.org/grpc/metadata"
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

type CommunityServer struct {
    community.UnimplementedCommunityServiceServer
    db *sql.DB
}

func New(dbR DbRepo, env string, rR RedisRepo, notCl NotificationS, cfg *config.Config) *Service {
	return &Service{
		dbR:   dbR,
		env:   env,
		rR:    rR,
		notCl: notCl,
	}
}

func NewCommunityServer(db *sql.DB) *CommunityServer {
    return &CommunityServer{
        db: db,
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

func (s *CommunityServer) InvitePeer(ctx context.Context, req *community.InvitePeerRequest) (*community.InvitePeerResponse, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "missing metadata")
    }

    initiatorUUIDs := md.Get("initiator_uuid")
    if len(initiatorUUIDs) == 0 {
        return nil, status.Error(codes.Unauthenticated, "initiator uuid not provided")
    }
    initiatorUUID := initiatorUUIDs[0]

    var initiatorID int
    var initiatorStatus string
    err := s.db.QueryRowContext(ctx,`SELECT id, status FROM participant WHERE link_edu = $1`, initiatorUUID,).Scan(&initiatorID, &initiatorStatus)
    if err == sql.ErrNoRows {
        return nil, status.Error(codes.PermissionDenied, "initiator not found")
    }
    if err != nil {
        return nil, status.Error(codes.Internal, "db error: "+err.Error())
    }
    if initiatorStatus != "ACTIVE" {
        return nil, status.Error(codes.PermissionDenied, "initiator not active")
    }

    var invitedID int
    err = s.db.QueryRowContext(ctx,`SELECT id FROM participant WHERE login = $1`, req.Login,
    ).Scan(&invitedID)
    if err == sql.ErrNoRows {
        return nil, status.Error(codes.NotFound, "invited login not found")
    }
    if err != nil {
        return nil, status.Error(codes.Internal, "db error: "+err.Error())
    }

    _, err = s.db.ExecContext(ctx,`INSERT INTO invites (initiator, invite_login) VALUES ($1, $2)`,
        initiatorUUID, req.Login,
    )
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to insert invite: "+err.Error())
    }

    return &community.InvitePeerResponse{Success: true}, nil
}
