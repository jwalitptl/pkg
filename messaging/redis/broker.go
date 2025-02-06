package redis

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"

	"aiclinic/pkg/messaging"
)

type redisBroker struct {
	client *redis.Client
}

func NewRedisBroker(redisURL string) (messaging.MessageBroker, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	return &redisBroker{client: client}, nil
}

func (b *redisBroker) Publish(ctx context.Context, channel string, message interface{}) error {
	msg := messaging.Message{
		Type:    channel,
		Payload: message,
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return b.client.Publish(ctx, channel, jsonMsg).Err()
}
