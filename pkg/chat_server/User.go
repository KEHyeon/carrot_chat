package chat_server

import (
	"carrot_chat/pkg/chat_protocol"
	"github.com/gorilla/websocket"
	"net"
)

// User 구조체는 각 사용자와 관련된 정보를 담습니다.
type User struct {
	Id   uint64          // 유저 ID
	Conn *websocket.Conn // WebSocket 연결
	Ip   net.IP          // 사용자 IP
}

func NewUser(id uint64, conn *websocket.Conn, ip net.IP) *User {
	return &User{
		Id:   id,
		Conn: conn,
		Ip:   ip,
	}
}

func (u *User) GetId() uint64 {
	return u.Id
}

// SendMessage는 사용자가 연결된 WebSocket을 통해 메시지를 전송합니다.
func (u *User) SendMessage(message chat_protocol.Message) error {
	// 메시지를 JSON으로 직렬화 후 전송
	err := u.Conn.WriteJSON(message)
	if err != nil {
		return err
	}
	return nil
}
