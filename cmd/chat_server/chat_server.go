package main

import (
	"carrot_chat/pkg/chat_server"
	"carrot_chat/pkg/config"
	"fmt"
)

func main() {
	cfg := config.NewConfig()
	// ChatServer 초기화
	server, err := chat_server.NewChatServer(cfg)
	if err != nil {
		panic(err)
	}
	// 서버 시작
	if err := server.Start(":8080"); err != nil {
		fmt.Println("Server Error:", err)
	}
}
