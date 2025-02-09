package main

import (
	"carrot_chat/pkg/chat_server"
	"fmt"
	"time"
)

func main() {
	// ChatServer 초기화
	server := chat_server.NewChatServer("your-secret-key", time.Hour)

	// 서버 시작
	if err := server.Start(":8080"); err != nil {
		fmt.Println("Server Error:", err)
	}
}
