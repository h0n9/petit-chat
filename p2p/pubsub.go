package p2p

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type (
	PubSub = pubsub.PubSub
	Topic  = pubsub.Topic
)

func (n *Node) NewPubSub() error {
	ps, err := pubsub.NewGossipSub(n.ctx, n.host)
	if err != nil {
		return nil
	}

	n.pubSub = ps

	return nil
}

func (n *Node) GetTopics() []string {
	return n.pubSub.GetTopics()
}

func (n *Node) SetTopic(topic string, data []byte) error {
	return n.pubSub.Publish(topic, data)
}

func (n *Node) Join(topic string) error {
	t, err := n.pubSub.Join(topic)
	if err != nil {
		return err
	}

	n.topics = append(n.topics, t)

	return nil
}
