package service

import (
	"context"
	community_proto "github.com/s21platform/community-proto/community-proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type Service struct {
	community_proto.UnimplementedCommunityServiceServer
	dbR DbRepo
}

func (s *Service) IsPeerExist(ctx context.Context, in *community_proto.EmailIn) (*community_proto.EmailOut, error) {
	log.Println("Input E-mail: ", in.Email)
	return &community_proto.EmailOut{IsExist: true}, nil
}

func (s *Service) SearchPeers(ctx context.Context, in *community_proto.SearchPeersIn) (*community_proto.SearchPeersOut, error) {
	log.Println("Input SearchPeers: ", in)
	res, err := s.dbR.SearchPeersBySubstring(ctx, in.Substring)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "search peer error: %s", err)
	}
	return &community_proto.SearchPeersOut{SearchPeers: res}, nil
}

func New(dbR DbRepo) *Service {
	return &Service{dbR: dbR}
}
