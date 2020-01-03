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

func (n *Node) GetSubs() []*Sub {
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
	sub, err := n.pubSub.Subscribe(topic)
	if err != nil {
		return err
	}

	n.subs = append(n.subs, sub)

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
	for i, s := range n.subs {
		if topic == s.Topic() {
			s.Cancel()
			n.subs = append(n.subs[:i], n.subs[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("not subscribed to '%s'", topic)
}
