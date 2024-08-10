package main

import (
	"fmt"
	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/repository/postgres"
)
import (
	"fmt"
	community_proto "github.com/s21platform/community-proto/community-proto"
	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	cfg := config.MustLoad()
	db := postgres.New(cfg)
	_ = db
	fmt.Println(cfg)
	cfg := config.MustLoad()

	thisService := service.New()

	s := grpc.NewServer()
	community_proto.RegisterCommunityServiceServer(s, thisService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("Cannnot listen port: %s; Error: %s", cfg.Service.Port, err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Cannnot start server: %s; Error: %s", cfg.Service.Port, err)
	}

}
