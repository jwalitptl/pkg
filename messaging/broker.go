package messaging

import (
	"context"
)

type MessageBroker interface {
	Publish(ctx context.Context, channel string, message interface{}) error
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
