package main

import (
	"fmt"
	"github.com/s21platform/community-service/internal/infra"
	"github.com/s21platform/metrics-lib/pkg"
	"log"
	"net"

	"google.golang.org/grpc"

	communityproto "github.com/s21platform/community-proto/community-proto"
	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/repository/postgres"
	"github.com/s21platform/community-service/internal/rpc"
)

func main() {
	cfg := config.MustLoad()
	dbRepo := postgres.New(cfg)

	thisService := rpc.New(dbRepo)

	metris, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "community", cfg.Platform.Env)
	if err != nil {
		log.Fatalf("cannot init metrics, err: %v", err)
	}
	defer metris.Disconnect()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.AuthInterceptor,
			infra.MetricsInterceptor(metris),
		),
	)
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
