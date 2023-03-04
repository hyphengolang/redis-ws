package pubsub

import (
	"context"

	"github.com/gobwas/ws/wsutil"
)

type PubSub interface {
	Publish(ctx context.Context, channel string, msg *wsutil.Message) error
	Subscribe(ctx context.Context, channel string) chan *wsutil.Message
}
