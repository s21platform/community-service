package service

import (
	"context"
	communityproto "github.com/s21platform/community-proto/community-proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type Service struct {
	communityproto.UnimplementedCommunityServiceServer
	dbR DbRepo
}

func (s *Service) IsPeerExist(ctx context.Context, in *communityproto.EmailIn) (*communityproto.EmailOut, error) {
	log.Println("Input E-mail: ", in.Email)
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
