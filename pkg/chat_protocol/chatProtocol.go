package chat_protocol

type Message struct {
	userId  uint64
	roomId  uint64
	message string
}
