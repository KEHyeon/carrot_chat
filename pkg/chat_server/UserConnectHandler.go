package chat_server

import (
	"carrot_chat/pkg/nats_client"
	redisclient "carrot_chat/pkg/redis_client"
	"carrot_chat/pkg/utils/jwtutil"
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type UserConnectHandler struct {
	jwtUtil     *jwtutil.JWTUtil // JWT 유틸리
	connections map[uint64]*User // 연결된 사용자 정보 관리
	mutex       sync.Mutex       // 동시 접근을 안전하게 처리하기 위한 뮤텍스
	redisClient *redisclient.RedisClient
	natsClient  *nats_client.NatsClient
}

func NewUserConnectHandler(jwtUtil *jwtutil.JWTUtil, redisClient *redisclient.RedisClient, natsClient *nats_client.NatsClient) *UserConnectHandler {
	return &UserConnectHandler{
		jwtUtil:     jwtUtil,
		redisClient: redisClient,
		connections: make(map[uint64]*User),
		mutex:       sync.Mutex{},
		natsClient:  natsClient,
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
	// 6-1. Redis에 현재 서버 정보 등록 (key = userID)
	serverIP := "현재 서버의 IP 주소" // 환경 변수에서 가져오거나 설정 파일에서 읽어올 것
	err = u.redisClient.Set(strconv.Itoa(int(userID)), serverIP)
	if err != nil {
		fmt.Printf("Redis 등록 실패: %v\n", err)
	}
	defer func() {
		//todo redis에서 삭제로직
		if err != nil {
			fmt.Printf("Redis 삭제 실패: %v\n", err)
		}
	}()
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

func (u *UserConnectHandler) getMessageHandler(user *User) {
	conn := user.Conn
	userID := user.GetId()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message for user %s: %v\n", userID, err)
			u.removeConnection(userID)
			break
		}

		fmt.Printf("Received message from user %s: %s\n", userID, len(message))
		err = u.handleMessage(message)
		if err != nil {
			fmt.Printf("Failed to publish message to NATS: %v\n", err)
		}
	}
}

// handleMessage는 수신된 메시지를 처리합니다.
func (u *UserConnectHandler) handleMessage(data []byte) error {
	return u.natsClient.PublishToQueue("chat", data)
}
