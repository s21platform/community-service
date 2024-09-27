package main

import (
	"fmt"
	communityproto "github.com/s21platform/community-proto/community-proto"
	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/repository/postgres"
	"github.com/s21platform/community-service/internal/service"
	"google.golang.org/grpc"
	"net"
)
import (
	"log"
)

func main() {
	cfg := config.MustLoad()
	dbRepo := postgres.New(cfg)

	thisService := service.New(dbRepo)

	s := grpc.NewServer()
	communityproto.RegisterCommunityServiceServer(s, thisService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("Cannnot listen port. Error: %s", err)
	}

	log.Println("Server is listening")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Cannnot start server. Error: %s", err)
	}

}
