package msg

import (
	"context"
	"time"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/types"
)

// Box refers to a chat room
type Box struct {
	ctx             context.Context
	topic           *types.Topic
	msgs            map[time.Time]*Msg
	latestTimestamp time.Time
}

func NewBox(ctx context.Context, topic *types.Topic) (*Box, error) {
	return &Box{
		ctx:             ctx,
		topic:           topic,
		msgs:            make(map[time.Time]*Msg),
		latestTimestamp: time.Now(),
	}, nil
}

func (b *Box) Publish(data []byte) error {
	if len(data) == 0 {
		// this is not error
		return nil
	}
	err := b.topic.Publish(b.ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (b *Box) Subscribe() (*types.Subscription, error) {
	return b.topic.Subscribe()
}

func (b *Box) Append(msg *Msg) error {
	_, exist := b.msgs[msg.Timestamp]
	if exist {
		return code.AlreadyAppendedMsg
	}

	if b.latestTimestamp.Before(msg.Timestamp) {
		b.latestTimestamp = msg.Timestamp
	}

	b.msgs[msg.Timestamp] = msg

	return nil
}
