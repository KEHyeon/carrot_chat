package chat_server

import (
	"carrot_chat/pkg/utils/jwtutil"
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"sync"
)

type UserConnectHandler struct {
	jwtUtil     *jwtutil.JWTUtil // JWT 유틸리티
	chatManager *ChatManager     // 채팅 매니저
	connections map[uint64]*User // 연결된 사용자 정보 관리
	mutex       sync.Mutex       // 동시 접근을 안전하게 처리하기 위한 뮤텍스
}

func NewUserConnectHandler(jwtUtil *jwtutil.JWTUtil, chatManager *ChatManager) *UserConnectHandler {
	return &UserConnectHandler{
		jwtUtil:     jwtUtil,
		chatManager: chatManager,
		connections: make(map[uint64]*User), // connections 맵 초기화
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 필요한 경우 Origin을 검증
		return true
	},
}

func (u *UserConnectHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 1. Query Parameter에서 JWT 토큰 추출
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Unauthorized: Token missing", http.StatusUnauthorized)
		return
	}

	// 2. JWT 토큰 검증
	claims, err := u.jwtUtil.ValidateToken(token)
	if err != nil {
		http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		return
	}

	// 3. 유저 정보 확인
	userID := claims.UserID
	fmt.Printf("User %s authenticated\n", userID)

	// 4. WebSocket 연결 업그레이드
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket Upgrade Error:", err)
		return
	}
	defer conn.Close()

	// 5. 연결된 유저 관리
	userIP := net.ParseIP(r.RemoteAddr) // 클라이언트 IP 추출
	user := NewUser(userID, conn, userIP)

	// 6. 유저를 connections 맵에 추가
	u.addUser(user)

	// 7. 메시지 처리 루프 시작
	u.getMessageHandler(user)
}

// addUser는 새로운 유저를 연결 목록에 안전하게 추가합니다.
func (u *UserConnectHandler) addUser(user *User) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	// 유저를 connections 맵에 추가
	u.connections[user.GetId()] = user
	fmt.Printf("User %d added to connections\n", user.GetId())
}

// removeConnection은 유저 연결을 안전하게 제거합니다.
func (u *UserConnectHandler) removeConnection(userID uint64) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if _, exists := u.connections[userID]; exists {
		fmt.Printf("Removing connection for room %d\n", userID)
		delete(u.connections, userID)
	}
}

// getMessageHandler는 유저의 메시지를 처리하는 루프입니다.
func (u *UserConnectHandler) getMessageHandler(user *User) {
	conn := user.Conn
	userID := user.GetId()

	// 메시지 처리 루프
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message for room %s: %v\n", userID, err)
			u.removeConnection(userID)
			break
		}

		// 클라이언트로부터 받은 메시지 처리
		fmt.Printf("Received message from room %s: %s\n", userID, string(message))

		// 메시지 처리 함수 호출
		handleMessage(user, message)
	}
}

// handleMessage는 수신된 메시지를 처리합니다.
func handleMessage(user *User, message []byte) {
	//message와 user를 바탕으로 nats에 넣을 메시지를 만든다.
	fmt.Println(string(message))
}
