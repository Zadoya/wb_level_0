package nats

// brew install nats-server
// чтобы установить натс сервер локально, но можно запускать из контейнера
// nats-server -- чтобы запустить натс сервер

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func Subscriber(ch chan []byte) error {
	nc, err := nats.Connect(nats.DefaultURL)//nats://localhost:4222
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	sub, err := nc.SubscribeSync("order")
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg, err := sub.NextMsg(time.Second * 1)
		if err != nil {
			if err.Error() == "nats: timeout" {
				continue
			}
			log.Fatal(err)
		}
		ch <- msg.Data
		log.Println("subscrider: received a message")
	}
}