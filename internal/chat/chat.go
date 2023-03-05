package chat

import (
	"context"
	"redis-ws/pkg/websocket"
	"redis-ws/pkg/websocket/pubsub"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Chat struct {
	ID string
}

// Broker is responsible of delegating the creation of a new Jam and the
// management of the Jam's websocket clients.
type Broker interface {
	Insert(ctx context.Context, id string) (*websocket.Client, bool)
}

type chatBroker struct {
	mu sync.Mutex
	m  map[string]*websocket.Client

	newClient func(ctx context.Context, id string) *websocket.Client
}

type Option func(*chatBroker)

func WithRedis(r *redis.Client) Option {
	return func(b *chatBroker) {
		b.newClient = func(ctx context.Context, id string) *websocket.Client {
			ps := pubsub.NewRedis(ctx, r, id)
			return websocket.NewClient(0, ps)
		}
	}
}

func NewBroker(opts ...Option) Broker {
	b := &chatBroker{
		m: make(map[string]*websocket.Client),
	}

	for _, opt := range opts {
		opt(b)
	}

	if b.newClient == nil {
		b.newClient = func(ctx context.Context, id string) *websocket.Client {
			ps := pubsub.New(256)
			return websocket.NewClient(0, ps)
		}
	}

	return b
}

func (b *chatBroker) Insert(ctx context.Context, id string) (*websocket.Client, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	found, ok := b.m[id]
	if ok {
		return found, ok
	}

	b.m[id] = b.newClient(ctx, id)
	return b.m[id], ok
}
