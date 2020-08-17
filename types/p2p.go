package types

import (
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type (
	Host   = host.Host
	PubSub = pubsub.PubSub
	Sub    = pubsub.Subscription
	Topic  = pubsub.Topic
)
