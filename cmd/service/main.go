package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/community-service/internal/client/notification"
	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/infra"
	"github.com/s21platform/community-service/internal/repository/postgres"
	"github.com/s21platform/community-service/internal/repository/redis"
	"github.com/s21platform/community-service/internal/service"
	"github.com/s21platform/community-service/pkg/community"
)

func main() {
	cfg := config.MustLoad()
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)

	dbRepo := postgres.New(cfg)
	defer dbRepo.Close()

	redisRepo := redis.New(cfg)

	notCl := notification.New(cfg)

	thisService := service.New(dbRepo, cfg.Platform.Env, redisRepo, notCl, cfg)

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, cfg.Service.Name, cfg.Platform.Env)
	if err != nil {
		log.Fatalf("cannot init metrics, err: %v", err)
	}
	defer metrics.Disconnect()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.Logger(logger),
			infra.AuthInterceptor,
			infra.MetricsInterceptor(metrics),
		),
	)
	community.RegisterCommunityServiceServer(s, thisService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("Cannnot listen port. Error: %s", err)
	}

	log.Println("Server is listening")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Cannnot start server. Error: %s", err)
	}
}
