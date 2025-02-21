package chat_protocol

type Message struct {
	fromUserId uint64
	roomId     uint64
	message    string
}
