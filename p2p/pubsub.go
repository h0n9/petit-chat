package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type PubSub = pubsub.PubSub

func NewPubSub(ctx context.Context, h host.Host) (*PubSub, error) {
	return pubsub.NewGossipSub(ctx, h)
}
