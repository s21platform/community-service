package main

import (
	"fmt"
	community_proto "github.com/s21platform/community-proto/community-proto"
	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/service"
	"google.golang.org/grpc"
	"net"
)
import (
	"log"
)

func main() {
	cfg := config.MustLoad()
	//_ = postgres.New(cfg)

	thisService := service.New()

	s := grpc.NewServer()
	community_proto.RegisterCommunityServiceServer(s, thisService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("Cannnot listen port. Error: %s", err)
	}

	log.Println("Server is listening")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Cannnot start server. Error: %s", err)
	}

}
