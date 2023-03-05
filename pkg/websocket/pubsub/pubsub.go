package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gobwas/ws/wsutil"
	"github.com/redis/go-redis/v9"
)

type PubSub interface {
	Publish(ctx context.Context, msg *wsutil.Message) error
	Subscribe(ctx context.Context) <-chan *wsutil.Message
}

type pubsub struct {
	broadcast chan *wsutil.Message
}

// Publish implements PubSub
func (ps *pubsub) Publish(ctx context.Context, msg *wsutil.Message) error {
	ps.broadcast <- msg
	return nil
}

// Subscribe implements PubSub
func (ps *pubsub) Subscribe(ctx context.Context) <-chan *wsutil.Message {
	return ps.broadcast
}

func New(cap int) PubSub {
	ps := &pubsub{
		broadcast: make(chan *wsutil.Message, cap),
	}
	return ps
}

type redisPS struct {
	r         *redis.Client
	channel   string
	broadcast chan *wsutil.Message
}

func NewRedis(ctx context.Context, r *redis.Client, channel string) PubSub {
	ps := &redisPS{
		r:         r,
		channel:   channel,
		broadcast: make(chan *wsutil.Message, 256),
	}
	go func() {
		for msg := range ps.r.Subscribe(ctx, ps.channel).Channel() {
			var payload wsutil.Message
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
				continue
			}
			ps.broadcast <- &payload
		}
	}()
	return ps
}

func (ps *redisPS) Publish(ctx context.Context, payload *wsutil.Message) error {
	p, err := json.Marshal(&payload)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return ps.r.Publish(context.Background(), ps.channel, string(p)).Err()
}

func (ps *redisPS) Subscribe(ctx context.Context) <-chan *wsutil.Message {
	return ps.broadcast
}
