package msg

import (
	"bytes"
	"context"
	"time"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/types"
)

const EOS = "EOS" // End Of Subscription

// Box refers to a chat room
type Box struct {
	ctx   context.Context
	myID  types.ID
	topic *types.Topic
	sub   *types.Sub

	msgs            map[time.Time]*Msg
	msgSubCh        chan *Msg
	latestTimestamp time.Time
}

func NewBox(ctx context.Context, myID types.ID, topic *types.Topic) (*Box, error) {
	return &Box{
		ctx:             ctx,
		myID:            myID,
		topic:           topic,
		sub:             nil,
		msgs:            make(map[time.Time]*Msg),
		msgSubCh:        nil,
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

func (b *Box) Subscribe() error {
	if b.sub != nil {
		return code.AlreadySubscribingTopic
	}

	sub, err := b.topic.Subscribe()
	if err != nil {
		return err
	}
	b.sub = sub

	for {
		received, err := sub.Next(b.ctx)
		if err != nil {
			return err
		}
		// TODO: consider if this a right way to handle closing subscription
		if bytes.Equal(received.GetData(), []byte(EOS)) {
			if received.GetFrom() == b.myID {
				sub.Cancel()
				err := b.topic.Close()
				if err != nil {
					return err
				}
				b.sub = nil
				break
			} else {
				continue
			}
		}
		msg := NewMsg(received.GetFrom(), received.GetData())
		err = b.append(msg)
		if err != nil {
			return err
		}
		if received.GetFrom() != b.myID && b.msgSubCh != nil {
			b.msgSubCh <- msg
		}
	}

	return nil
}

func (b *Box) Close() error {
	return b.Publish([]byte(EOS))
}

func (b *Box) Subscribing() bool {
	return b.sub != nil
}

func (b *Box) SetMsgSubCh(msgSubCh chan *Msg) {
	b.msgSubCh = msgSubCh
}

func (b *Box) GetMsgs() map[time.Time]*Msg {
	return b.msgs
}

func (b *Box) append(msg *Msg) error {
	timestamp := msg.GetTime()
	_, exist := b.msgs[timestamp]
	if exist {
		return code.AlreadyAppendedMsg
	}

	if b.latestTimestamp.Before(timestamp) {
		b.latestTimestamp = timestamp
	}

	b.msgs[timestamp] = msg

	return nil
}
