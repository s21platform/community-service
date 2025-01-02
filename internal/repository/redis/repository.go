package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func (r *Repository) Get(ctx context.Context) (string, error) {
	keys, err := r.conn.Keys(ctx, "*").Result()
	if err != nil {
		fmt.Printf("cannot get keys, err: %v", err)
		return "", fmt.Errorf("cannot get keys, err: %v", err)
	}

	if len(keys) > 0 {
		randomKey := keys[rand.Intn(len(keys))]

		val, err := r.conn.Get(ctx, randomKey).Result()
		if errors.Is(err, redis.Nil) {
			log.Printf("cannot find key %s \n", randomKey)
			return "", err
		} else if err != nil {
			log.Printf("cannot get value by key: %s, err: %v\n", randomKey, err)
			return "", fmt.Errorf("cannot get value by key: %s, err: %v", randomKey, err)
		}
		return val, nil
	} else {
		return "", status.Errorf(codes.Unknown, "no key in Redis")
	}
}

func (r *Repository) GetByKey(ctx context.Context, key string) (string, error) {
	val, err := r.conn.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", fmt.Errorf("cannot get value by key: %s, err: %v", key, err)

	}
	return val, nil
}

func (r *Repository) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := r.conn.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Printf("cannot set key: %s, value: %s, err: %v", key, value, err)
		return fmt.Errorf("cannot set key: %s, value: %s, err: %v", key, value, err)
	}
	return nil
}
