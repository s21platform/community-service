package main

import "github.com/s21platform/community-service/internal/config"

func main() {
	_ = config.MustLoad()
}
