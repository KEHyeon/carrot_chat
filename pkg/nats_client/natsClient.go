package nats_client

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

// natsClient 구조체 정의
type natsClient struct {
	conn *nats.Conn
}

// newNatsClient 함수 구현
func newNatsClient(url string) (*natsClient, error) {
	// NATS 서버에 연결 시도
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %v", err)
	}

	// 성공적으로 연결되면 natsClient 객체 반환
	client := &natsClient{
		conn: conn,
	}

	return client, nil
}

// NATS 클라이언트 연결 해제 함수
func (nc *natsClient) close() {
	if nc.conn != nil {
		nc.conn.Close()
		log.Println("NATS connection closed.")
	}
}
