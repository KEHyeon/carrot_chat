package chat_server

import (
	"carrot_chat/pkg/chat_protocol"
	"sync"
	"sync/atomic"
)

type ChatManager struct {
	ChatMap map[uint64]*Chat
	chatId  atomic.Uint64
	mutex   sync.Mutex // Mutex 추가
}

// NewChatManager는 새로운 ChatManager를 초기화하는 함수입니다.
func NewChatManager() *ChatManager {
	return &ChatManager{
		ChatMap: make(map[uint64]*Chat), // RoomMap을 빈 맵으로 초기화
		chatId:  atomic.Uint64{},        // roomId는 기본값인 0으로 초기화
	}
}

// AddRoom은 채팅방을 ChatManager에 추가합니다.
func (c *ChatManager) AddRoom(room *Chat) {
	c.mutex.Lock()         // 동기화 시작
	defer c.mutex.Unlock() // 동기화 끝

	nextRoomId := c.chatId.Add(1)
	c.ChatMap[nextRoomId] = room
}

// MessageHandler는 메시지를 처리하는 함수입니다.
func (c *ChatManager) MessageHandler(message chat_protocol.Message) {
	// 메시지 처리 로직
	c.mutex.Lock()         // 동기화 시작
	defer c.mutex.Unlock() // 동기화 끝

	// ChatMap 수정 로직이 필요하면 여기에서 처리
}
