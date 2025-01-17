package rpc

import (
	"context"
	"log"

	communityproto "github.com/s21platform/community-proto/community-proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	communityproto.UnimplementedCommunityServiceServer
	dbR DbRepo
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
	peerStatus, err := s.dbR.GetPeerStatus(ctx, in.Email)
	if err != nil {
		log.Printf("cannot get peer status, err: %v\n", err)
		return nil, status.Errorf(codes.Internal, "cannot check peer error: %s", err)
	}

	if peerStatus != "ACTIVE" {
		log.Printf("peer status: %s\n", peerStatus)
		return &communityproto.EmailOut{IsExist: false}, nil
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

func New(dbR DbRepo) *Service {
	return &Service{dbR: dbR}
}
