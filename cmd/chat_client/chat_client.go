package main

import (
	"fmt"
	"log"
	"time"

	"carrot_chat/pkg/utils/jwtutil"
	"github.com/gorilla/websocket"
)

func main() {
	// JWT 유틸리티 초기화
	jwtUtil := jwtutil.NewJWTUtil("your-secret-key", 15*24*time.Hour) // 만료 시간을 하루로 설정

	// 유저 ID를 예시로 사용하여 JWT 토큰 생성
	userID := uint64(12345)
	token, err := jwtUtil.GenerateToken(userID)
	if err != nil {
		log.Fatal("Error generating token:", err)
	}

	// WebSocket 서버 주소와 JWT 토큰을 포함하여 URL 생성
	serverURL := fmt.Sprintf("ws://localhost:8080/ws?token=%s", token)

	// WebSocket 연결 시도
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer conn.Close()

	fmt.Println("Connected to WebSocket server")

	// 서버와 메시지 송수신
	message := "Hello, Server!"
	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Fatal("Error sending message:", err)
	}
	fmt.Printf("Sent message: %s\n", message)

	// 서버로부터 응답 받기
	_, response, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("Error reading message:", err)
	}
	fmt.Printf("Received response: %s\n", string(response))

	// 종료
	fmt.Println("Closing connection")
}
