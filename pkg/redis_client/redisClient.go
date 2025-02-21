package redisclient

import (
	"carrot_chat/pkg/config"
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(cfg *config.Config) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis 연결 실패: %v", err)
	} else {
		fmt.Println("Redis에 연결되었습니다.")
	}

	return &RedisClient{client: client}
}

// Set은 Redis에 키-값을 저장합니다.
func (r *RedisClient) Set(key string, value interface{}) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

// Get은 Redis에서 키의 값을 조회합니다.
func (r *RedisClient) Get(key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// Close는 Redis 클라이언트를 종료합니다.
func (r *RedisClient) Close() {
	r.client.Close()
}
