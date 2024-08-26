package service

import (
	"context"
	community_proto "github.com/s21platform/community-proto/community-proto"
	"log"
)

type Service struct {
	community_proto.UnimplementedCommunityServiceServer
}

func (s *Service) IsPeerExist(ctx context.Context, in *community_proto.EmailIn) (*community_proto.EmailOut, error) {
	log.Println("Input E-mail: ", in.Email)
	return &community_proto.EmailOut{IsExist: true}, nil
}

func New() *Service {
	return &Service{}
}
