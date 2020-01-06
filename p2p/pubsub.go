package p2p

import (
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/h0n9/petit-chat/code"
)

type (
	Topic = pubsub.Topic
	Sub   = pubsub.Subscription
)

func (n *Node) NewPubSub() error {
	ps, err := pubsub.NewGossipSub(n.ctx, n.host)
	if err != nil {
		return err
	}

	n.pubSub = ps

	return nil
}

func (n *Node) publish(topic string, data []byte) error {
	err := n.pubSub.Publish(topic, data)
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) subscribe(topic string) error {
	_, exist := n.subs[topic]
	if exist {
		return code.NonSubscribingTopic
	}

	sub, err := n.pubSub.Subscribe(topic)
	if err != nil {
		return err
	}

	n.subs[topic] = sub

	go func() {
		for {
			msg, err := sub.Next(n.ctx)
			if err != nil {
				continue
			}
			if len(msg.Data) == 0 {
				continue
			}
			if msg.GetFrom().String() == n.host.ID().String() {
				continue
			}

			fmt.Printf("\x1b[32m%s: %s\x1b[0m\n> ",
				msg.GetFrom(),
				msg.GetData(),
			)
		}
	}()

	return nil
}

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
