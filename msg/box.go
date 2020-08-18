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
	topic *types.Topic
	sub   *types.Sub

	myID            types.ID
	msgSubCh        chan *Msg
	latestTimestamp time.Time
	readUntilIndex  int

	msgs      []*Msg
	msgHashes map[types.Hash]*Msg
}

func NewBox(ctx context.Context, topic *types.Topic, myID types.ID) (*Box, error) {
	return &Box{
		ctx:   ctx,
		topic: topic,
		sub:   nil,

		myID:            myID,
		msgSubCh:        nil,
		latestTimestamp: time.Now(),
		readUntilIndex:  0,

		msgs:      make([]*Msg, 0),
		msgHashes: make(map[types.Hash]*Msg),
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
		readUntilIndex, err := b.append(msg)
		if err != nil {
			return err
		}
		if received.GetFrom() != b.myID && b.msgSubCh != nil {
			b.msgSubCh <- msg
			b.readUntilIndex = readUntilIndex
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

func (b *Box) GetMsgs() []*Msg {
	return b.msgs
}

func (b *Box) GetMsgsFromRead() []*Msg {
	return append([]*Msg{}, b.msgs[b.readUntilIndex:]...)
}

func (b *Box) append(msg *Msg) (int, error) {
	hash, err := msg.Hash()
	if err != nil {
		return 0, err
	}

	_, exist := b.msgHashes[hash]
	if exist {
		return 0, code.AlreadyAppendedMsg
	}

	timestamp := msg.GetTime()
	if b.latestTimestamp.Before(timestamp) {
		b.latestTimestamp = timestamp
	}

	b.msgs = append(b.msgs, msg)
	b.msgHashes[hash] = msg

	return len(b.msgs) - 1, nil
}
