package nats_client

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

// natsClient 구조체 정의
type NatsClient struct {
	conn *nats.Conn
}

// newNatsClient 함수 구현
func NewNatsClient(url string) (*NatsClient, error) {
	// NATS 서버에 연결 시도
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %v", err)
	}

	// 성공적으로 연결되면 natsClient 객체 반환
	client := &NatsClient{
		conn: conn,
	}

	return client, nil
}

// NATS 클라이언트 연결 해제 함수
func (nc *NatsClient) close() {
	if nc.conn != nil {
		nc.conn.Close()
		log.Println("NATS connection closed.")
	}
}

// 메시지를 큐에 넣는 함수
func (nc *NatsClient) PublishToQueue(subject string, message []byte) error {
	err := nc.conn.Publish(subject, message)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	log.Printf("Message published to subject '%s': %s\n", subject, message)
	return nil
}

// 큐에서 메시지를 읽는 함수
func (nc *NatsClient) SubscribeFromQueue(subject, queue string) error {
	_, err := nc.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		// msg.Data는 []byte 타입
		log.Printf("Received message: %v\n", msg.Data) // []byte 그대로 출력
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to queue: %v", err)
	}
	log.Printf("Subscribed to subject '%s' on queue '%s'.\n", subject, queue)
	return nil
}
