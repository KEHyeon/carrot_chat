package chat_server

import (
	"carrot_chat/pkg/config"
	natsClient "carrot_chat/pkg/nats_client"
	redisclient "carrot_chat/pkg/redis_client"
	"carrot_chat/pkg/utils/jwtutil"
	"fmt"
	"net/http"
)

// ChatServer는 WebSocket 서버와 관련된 설정과 시작 로직을 포함하는 구조체입니다.
type ChatServer struct {
	jwtUtil            *jwtutil.JWTUtil
	userConnectHandler *UserConnectHandler
}

func NewChatServer(cfg *config.Config) (*ChatServer, error) {
	jwtUtil := jwtutil.NewJWTUtil(cfg.SecretKey, cfg.TokenExpireDuration)
	redisClient, err := redisclient.NewRedisClient(cfg)
	if err != nil {
		return nil, err
	}
	natsClient, err := natsClient.NewNatsClient(cfg.NatsUrl)
	if err != nil {
		return nil, err
	}
	userConnectHandler := NewUserConnectHandler(jwtUtil, redisClient, natsClient)

	return &ChatServer{
		jwtUtil:            jwtUtil,
		userConnectHandler: userConnectHandler,
	}, nil
}

// Start는 서버를 시작하고 WebSocket 핸들러를 등록합니다.
func (s *ChatServer) Start(address string) error {
	// WebSocket 핸들러 등록
	http.HandleFunc("/ws", s.userConnectHandler.HandleWebSocket)

	// 서버 시작
	fmt.Printf("Starting server on %s...\n", address)
	return http.ListenAndServe(address, nil)
}
