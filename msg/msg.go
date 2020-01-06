package msg

import (
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type (
	PubSub = pubsub.PubSub
	Sub    = pubsub.Subscription
)

type Msg struct {
	Timestamp time.Time

	Incoming bool
	Read     bool
	Received bool

	Value pubsub.Message
}

func NewMsg(value pubsub.Message) *Msg {
	return &Msg{
		Timestamp: time.Now(),
		Read:      false,
		Received:  true,
		Value:     value,
	}
}
