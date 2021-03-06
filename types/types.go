package types

import (
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type (
	ID        = peer.ID
	Host      = host.Host
	PubSub    = pubsub.PubSub
	Sub       = pubsub.Subscription
	Topic     = pubsub.Topic
	PubSubMsg = pubsub.Message
	Hash      = [32]byte
)
