package chat_server

import (
	"carrot_chat/pkg/utils/jwtutil"
	"fmt"
	"net/http"
	"time"
)

// ChatServer는 WebSocket 서버와 관련된 설정과 시작 로직을 포함하는 구조체입니다.
type ChatServer struct {
	jwtUtil            *jwtutil.JWTUtil // JWT 유틸리티
	chatManager        *ChatManager
	userConnectHandler *UserConnectHandler
}

// NewChatServer는 ChatServer 구조체를 초기화하고 반환합니다.
func NewChatServer(secretKey string, tokenExpireDuration time.Duration) *ChatServer {
	jwtUtil := jwtutil.NewJWTUtil(secretKey, tokenExpireDuration) // JWT 유틸리티 초기화
	chatManager := NewChatManager()                               // 채팅 매니저 초기화
	userConnectHandler := NewUserConnectHandler(jwtUtil, chatManager)

	return &ChatServer{
		jwtUtil:            jwtUtil,
		chatManager:        chatManager,
		userConnectHandler: userConnectHandler,
	}
}

// Start는 서버를 시작하고 WebSocket 핸들러를 등록합니다.
func (s *ChatServer) Start(address string) error {
	// WebSocket 핸들러 등록
	http.HandleFunc("/ws", s.userConnectHandler.HandleWebSocket)

	// 서버 시작
	fmt.Printf("Starting server on %s...\n", address)
	return http.ListenAndServe(address, nil)
}
