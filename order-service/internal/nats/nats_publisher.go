package nats

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type NatsPublisher struct {
	conn *nats.Conn
}

func NewNatsPublisher(url string) *NatsPublisher {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	return &NatsPublisher{conn: nc}
}

func (p *NatsPublisher) PublishOrderCreated(subject string, order interface{}) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return p.conn.Publish(subject, data)
}
