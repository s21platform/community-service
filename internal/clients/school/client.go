package school

import (
	"context"
	"fmt"
	"log"

	"github.com/s21platform/community-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	schoolproto "github.com/s21platform/school-proto/school-proto"
)

// Создаем обертку над клиентом 
type Handle struct {
	client schoolproto.SchoolServiceClient
}

// метод, который дергает ручку скул сервиса
func (h *Handle) GetCampuses(ctx context.Context) () {}

// Создаем и настраиваем соединение с gRPC-сервисом School
func MustConnect(cfg *config.Config) *Handle {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.School.Host, cfg.School.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to community service: %v", err)
	}
	client := schoolproto.NewSchoolServiceClient(conn)
	return &Handle{client: client}
}