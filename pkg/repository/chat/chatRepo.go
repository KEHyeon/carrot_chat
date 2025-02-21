package chat

import (
	"carrot_chat/pkg/chat_protocol"
	"carrot_chat/pkg/chat_server"
)

type chat interface {
	SendMessage(messege chat_protocol.Message) (*chat_server.Ch, error)
}
