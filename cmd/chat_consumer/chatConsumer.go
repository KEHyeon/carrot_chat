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

	// 구독 설정
	sub, err := nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("수신된 메시지: %s\n", string(m.Data))
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	// 메시지 발행
	err = nc.Publish("foo", []byte("Hello, NATS!"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("메시지가 발행되었습니다.")

	// 잠시 대기하여 메시지를 받을 수 있도록 함
	select {}
}
