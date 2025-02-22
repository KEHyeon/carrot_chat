package main

import (
	"bufio"
	"carrot_chat/pkg/chat/pb/proto"
	"carrot_chat/pkg/utils/jwtutil"
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	// JWT 유틸리티 초기화
	jwtUtil := jwtutil.NewJWTUtil("your-secret-key", 15*24*time.Hour) // 만료 시간을 하루로 설정
	userID := uint64(12345)

	// JWT 토큰 생성
	token, err := jwtUtil.GenerateToken(userID)
	if err != nil {
		log.Fatal("Error generating token:", err)
	}

	// WebSocket 서버 URL 생성
	serverURL := fmt.Sprintf("ws://localhost:8080/ws?token=%s", token)

	// WebSocket 연결 시도
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer conn.Close()

	fmt.Println("✅ Connected to WebSocket server")

	// 채널을 사용하여 종료 신호를 감지
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// 메시지 송수신을 위한 고루틴 실행
	go writePump(conn)
	go readPump(conn)

	// 프로그램 종료 대기
	<-interrupt
	fmt.Println("🚪 Closing connection...")
	conn.Close()
}
func writePump(conn *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin) // CMD에서 입력 받기 위한 Scanner 설정
	fmt.Println("채팅 메시지를 입력하세요. (종료하려면 'exit' 입력)")

	for {
		fmt.Print("채팅방 ID: ")
		scanner.Scan()
		roomID := scanner.Text()

		if roomID == "exit" {
			fmt.Println("프로그램을 종료합니다.")
			return
		}

		fmt.Print("메시지 내용: ")
		scanner.Scan()
		content := scanner.Text()

		fmt.Print("메시지 타입 (TEXT/IMAGE/FILE 등): ")
		scanner.Scan()
		messageTypeStr := scanner.Text()
		var messageType chatpb.MessageType
		switch messageTypeStr {
		case "TEXT":
			messageType = chatpb.MessageType_TEXT
		case "IMAGE":
			messageType = chatpb.MessageType_IMAGE
		case "FILE":
			messageType = chatpb.MessageType_VIDEO
		default:
			messageType = chatpb.MessageType_TEXT // 기본값 TEXT
		}

		// 현재 시간 (Unix Timestamp)
		timestamp := time.Now().Unix()

		// 채팅 메시지 생성
		chatMessage := &chatpb.ChatMessage{
			UserId:    12345, // 예시로 유저 ID
			RoomId:    roomID,
			Content:   content,
			Type:      messageType,
			Timestamp: timestamp,
		}

		// Protobuf 메시지 직렬화
		data, err := proto.Marshal(chatMessage)
		if err != nil {
			log.Println("❌ Error marshalling message:", err)
			return
		}

		// 메시지 전송
		err = conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			log.Println("❌ Error sending message:", err)
			return
		}
		fmt.Printf("📤 Sent message: %+v\n", chatMessage)
	}
}

// 서버로부터 메시지를 비동기적으로 수신하는 함수
func readPump(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ Error reading message:", err)
			return
		}
		fmt.Printf("📥 Received: %s\n", string(message))
	}
}
