package chat_server

import (
	"carrot_chat/pkg/chat_protocol"
	"sync/atomic"
)

type ChatManager struct {
	RoomMap map[uint64]*Room
	roomId  atomic.Uint64
}

// NewChatManager는 새로운 ChatManager를 초기화하는 함수입니다.
func NewChatManager() *ChatManager {
	return &ChatManager{
		RoomMap: make(map[uint64]*Room), // RoomMap을 빈 맵으로 초기화
		roomId:  atomic.Uint64{},        // roomId는 기본값인 0으로 초기화
	}
}
func (c *ChatManager) AddRoom(room *Room) {
	nextRoomId := c.roomId.Add(1)
	c.RoomMap[nextRoomId] = room
}

func (c *ChatManager) MessageHandler(message chat_protocol.Message) {

}
