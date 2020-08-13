package p2p

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type (
	PubSub = pubsub.PubSub
	Sub    = pubsub.Subscription
)

func (n *Node) NewPubSub() error {
	ps, err := pubsub.NewGossipSub(n.ctx, n.host)
	if err != nil {
		return err
	}

	n.pubSub = ps

	return nil
}

/*
func (n *Node) unsubscribe(topic string) error {
	for t, sub := range n.subs {
		if topic == t {
			sub.Cancel()
			delete(n.subs, t)

			return nil
		}
	}

	return code.NonSubscribingTopic
}
*/
