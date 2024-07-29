package main

import (
	"fmt"
	"github.com/s21platform/community-service/internal/config"
	"github.com/s21platform/community-service/internal/repository/postgres"
)

func main() {

	cfg := config.MustLoad()
	db := postgres.New(cfg)
	_ = db
	fmt.Println(cfg)
}
