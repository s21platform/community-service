package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/s21platform/community-service/internal/config"
)

type Repository struct {
	conn *redis.Client
}

func New(cfg *config.Config) *Repository {
	redisPort := cfg.Cache.Port
	redisHost := cfg.Cache.Host
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	rdb := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     "",
		DB:           0,
		MinIdleConns: 2,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	return &Repository{conn: rdb}
}

func (r *Repository) GetByKey(ctx context.Context, key config.Key) (string, error) {
	val, err := r.conn.Get(ctx, string(key)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", fmt.Errorf("cannot get value by key: %s, err: %v", key, err)
	}
	return val, nil
}

func (r *Repository) Set(ctx context.Context, key config.Key, value string, expiration time.Duration) error {
	err := r.conn.Set(ctx, string(key), value, expiration).Err()
	if err != nil {
		log.Printf("cannot set key: %s, value: %s, err: %v", key, value, err)
		return fmt.Errorf("cannot set key: %s, value: %s, err: %v", key, value, err)
	}
	return nil
}
