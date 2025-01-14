package main

import (
	"fmt"
	logger_lib "github.com/s21platform/logger-lib"
	"log"
	"net"

	communityproto "github.com/s21platform/community-proto/community-proto"
	"github.com/s21platform/metrics-lib/pkg"
	"google.golang.org/grpc"

	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/infra"
	"github.com/s21platform/community-service/internal/repository/postgres"
	"github.com/s21platform/community-service/internal/rpc"
)

func main() {
	cfg := config.MustLoad()
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)

	dbRepo := postgres.New(cfg)

	thisService := rpc.New(dbRepo)

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "community", cfg.Platform.Env)
	if err != nil {
		log.Fatalf("cannot init metrics, err: %v", err)
	}
	defer metrics.Disconnect()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.Logger(logger),
			infra.AuthInterceptor(cfg.Platform.Env, *dbRepo),
			infra.MetricsInterceptor(metrics),
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
