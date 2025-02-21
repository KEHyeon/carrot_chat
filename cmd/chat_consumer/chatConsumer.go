package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	// NATS 서버에 연결
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// 메시지 큐 구독
	// "foo" 큐에 메시지를 구독하여 처리하는 예제
	sub, err := nc.QueueSubscribe("foo", "worker-group", func(m *nats.Msg) {
		fmt.Printf("수신된 메시지: %s\n", string(m.Data))
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	// 메시지 큐에 푸시하는 발행자
	// 발행자는 "foo" 채널로 메시지를 푸시
	err = nc.Publish("foo", []byte("Hello, NATS Queue!"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("메시지가 큐에 푸시되었습니다.")

	// 잠시 대기하여 메시지를 받을 수 있도록 함
	select {}
}
