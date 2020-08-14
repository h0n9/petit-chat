package p2p

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/h0n9/petit-chat/types"
)

func (n *Node) NewPubSub() error {
	ps, err := pubsub.NewGossipSub(n.ctx, n.host)
	if err != nil {
		return err
	}

	n.pubSub = ps

	return nil
}

func (n *Node) Join(topic string) (*types.Topic, error) {
	tp, err := n.pubSub.Join(topic)
	if err != nil {
		return nil, err
	}
	return tp, nil
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
