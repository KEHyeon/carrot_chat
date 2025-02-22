package config

import (
	"github.com/nats-io/nats.go"
	"time"
)

type Config struct {
	SecretKey           string
	TokenExpireDuration time.Duration
	RedisAddr           string
	RedisPassword       string
	RedisDB             int
	Port                int
	NatsUrl             string
}

// NewConfig는 기본 설정을 반환합니다.
func NewConfig() *Config {
	return &Config{
		SecretKey:           "your-secret-key",
		TokenExpireDuration: time.Hour,
		RedisAddr:           "localhost:6379",
		RedisPassword:       "",
		RedisDB:             0,
		Port:                8081,
		NatsUrl:             nats.DefaultURL,
	}
}
