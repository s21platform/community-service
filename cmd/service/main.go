package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/community-service/internal/api"
	"github.com/s21platform/community-service/internal/client/notification"
	"github.com/s21platform/community-service/internal/config"
	apigen "github.com/s21platform/community-service/internal/generated"
	"github.com/s21platform/community-service/internal/infra"
	"github.com/s21platform/community-service/internal/repository/postgres"
	"github.com/s21platform/community-service/internal/repository/redis"
	"github.com/s21platform/community-service/internal/service"
	"github.com/s21platform/community-service/pkg/community"
)

func main() {
	cfg := config.MustLoad()
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)
	ctx := logger_lib.NewContext(context.Background(), logger)

	dbRepo := postgres.New(cfg)
	defer dbRepo.Close()

	redisRepo := redis.New(cfg)

	notificationClient := notification.New(cfg)

	thisService := service.New(dbRepo, cfg.Platform.Env, redisRepo, notificationClient, cfg)

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, cfg.Service.Name, cfg.Platform.Env)
	if err != nil {
		log.Fatalf("cannot init metrics, err: %v", err)
	}
	defer metrics.Disconnect()

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.LoggerRPC(logger),
			infra.AuthInterceptor,
			infra.MetricsInterceptor(metrics),
		),
	)
	community.RegisterCommunityServiceServer(grpcSrv, thisService)

	handler := api.New(dbRepo, redisRepo, notificationClient)
	router := chi.NewRouter()
	router.Use(infra.AuthRequest)
	router.Use(infra.LoggerHTTP(logger))

	apigen.HandlerFromMux(handler, router)
	httpServer := &http.Server{
		Handler: router,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("failed to listen port; error: %s", err))
		os.Exit(1)
	}

	m := cmux.New(lis)

	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	g, _ := errgroup.WithContext(ctx)

	logger_lib.Info(ctx, "starting server")

	g.Go(func() error {
		if err := grpcSrv.Serve(grpcListener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return fmt.Errorf("gRPC server error: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := httpServer.Serve(httpListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server error: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := m.Serve(); err != nil {
			return fmt.Errorf("cannot start service: %v", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("server error: %v", err))
		os.Exit(1)
	}
}
