package msg

import (
	"context"
	"fmt"
	"time"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/types"
)

// Box refers to a chat room
type Box struct {
	ctx   context.Context
	topic *types.Topic
	sub   *types.Sub

	myID            types.ID
	msgSubCh        chan *Msg
	latestTimestamp time.Time
	readUntilIndex  int

	msgs      []*Msg              // TODO: limit the size of msgs slice
	msgHashes map[types.Hash]*Msg // TODO: limit the size of msgHashes map
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

func (b *Box) Publish(t types.Msg, parentMsgHash types.Hash, data []byte) error {
	if len(data) == 0 {
		// this is not error
		return nil
	}
	msg := NewMsg(b.myID, t, parentMsgHash, data)
	data, err := msg.Encapsulate()
	if err != nil {
		return err
	}
	err = b.topic.Publish(b.ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (b *Box) Subscribe(handler MsgHandler) error {
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
			// TODO: replace fmt.Println() to logger.Println()
			fmt.Println(err)
			continue
		}
		eos, err := handler(b, received)
		if err != nil {
			// TODO: replace fmt.Println() to logger.Println()
			fmt.Println(err)
			continue
		}

		// eos shoud be the only way to break for loop
		if eos {
			b.sub.Cancel()
			err = b.topic.Close()
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}

	return nil
}

func (b *Box) Close() error {
	// Announe EOS to others (application layer)
	return b.Publish(types.MsgEOS, types.Hash{}, []byte("bye"))
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

func (b *Box) GetMsg(mh types.Hash) *Msg {
	return b.msgHashes[mh]
}

func (b *Box) GetUnreadMsgs() []*Msg {
	msgs := []*Msg{}
	if b.readUntilIndex+1 < len(b.msgs) {
		msgs = append(msgs, b.msgs[b.readUntilIndex+1:]...)
	}
	b.readUntilIndex = len(b.msgs) - 1
	return msgs
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
