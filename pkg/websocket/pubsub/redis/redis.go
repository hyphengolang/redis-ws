package pubsub

import (
	"context"
	"encoding/json"
	"redis-ws/pkg/websocket/pubsub"

	"github.com/gobwas/ws/wsutil"
	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	r  *redis.Client
	ch chan *wsutil.Message
}

func New(r *redis.Client, cap int) pubsub.PubSub {
	ps := redisClient{
		r:  r,
		ch: make(chan *wsutil.Message, cap),
	}

	return &ps
}

func (ps *redisClient) Publish(ctx context.Context, channel string, msg *wsutil.Message) error {
	p, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return ps.r.Publish(ctx, channel, string(p)).Err()
}

func (ps *redisClient) Subscribe(ctx context.Context, channel string) chan *wsutil.Message {
	go func() {
		// TODO -- buffer connection
		_msg := <-ps.r.Subscribe(ctx, channel).Channel()

		var msg wsutil.Message
		_ = json.Unmarshal([]byte(_msg.Payload), &msg)

		ps.ch <- &msg
	}()

	return ps.ch
}
