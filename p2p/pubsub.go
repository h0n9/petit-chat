package p2p

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type (
	PubSub = pubsub.PubSub
	Topic  = pubsub.Topic
	Sub    = pubsub.Subscription
)

func (n *Node) NewPubSub() error {
	ps, err := pubsub.NewGossipSub(n.ctx, n.host)
	if err != nil {
		return nil
	}

	n.pubSub = ps

	return nil
}

func (n *Node) GetSubs() map[string]*Sub {
	return n.subs
}

func (n *Node) GetPeers(topic string) []peer.ID {
	return n.pubSub.ListPeers(topic)
}

func (n *Node) Publish(topic string, data []byte) error {
	err := n.pubSub.Publish(topic, data)
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) Subscribe(topic string) error {
	_, exist := n.subs[topic]
	if exist {
		return fmt.Errorf("already subscried to %s", topic)
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

func (n *Node) Unsubscribe(topic string) error {
	for t, sub := range n.subs {
		if topic == t {
			sub.Cancel()
			delete(n.subs, t)

			return nil
		}
	}

	return fmt.Errorf("not subscribed to '%s'", topic)
}
