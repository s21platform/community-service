package service

import (
	"context"
	community_proto "github.com/s21platform/community-proto/community-proto"
)

type Service struct {
	community_proto.UnimplementedCommunityServiceServer
}

func (s *Service) IsPeerExist(ctx context.Context, eIn *community_proto.EmailIn) (*community_proto.EmailOut, error) {
	return &community_proto.EmailOut{IsExist: true}, nil
}

func New() *Service {
	return &Service{}
}
