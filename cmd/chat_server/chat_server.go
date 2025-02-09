package main

import (
	"carrot_chat/pkg/chat_server"
	"carrot_chat/pkg/utils/jwtutil"
	"fmt"
	"net/http"
	"time"
)

func main() {
	// JWT 유틸리티와 채팅 매니저 초기화
	jwtUtil := jwtutil.NewJWTUtil("your-secret-key", time.Hour) // 여기에 맞는 JWT 유틸리티 초기화
	chatManager := chat_server.NewChatManager()                 // 채팅 매니저 초기화

	// UserConnectHandler 생성
	userConnectHandler := chat_server.NewUserConnectHandler(jwtUtil, chatManager)

	// WebSocket 핸들러 등록
	http.HandleFunc("/ws", userConnectHandler.HandleWebSocket) // 여기서 userConnectHandler의 메서드 사용

	// 서버 시작
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server Error:", err)
	}
}
