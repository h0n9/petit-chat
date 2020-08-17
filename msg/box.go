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
	ctx     context.Context
	topic   *types.Topic
	sub     *types.Sub
	closeCh chan bool

	msgs            map[time.Time]*Msg
	latestTimestamp time.Time
}

func NewBox(ctx context.Context, topic *types.Topic) (*Box, error) {
	return &Box{
		ctx:             ctx,
		topic:           topic,
		sub:             nil,
		closeCh:         make(chan bool, 1),
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

func (b *Box) Open() error {
	if b.sub != nil {
		return code.AlreadySubscribingTopic
	}

	sub, err := b.topic.Subscribe()
	if err != nil {
		return err
	}
	b.sub = sub

	fmt.Println("open")

	errs := make(chan error, 1)

	// routine for receiving data
	go func() {
		fmt.Println("start routine")
		for {
			select {
			case <-b.closeCh:
				sub.Cancel()
				b.sub = nil
				break
			default:
			}
			received, err := sub.Next(b.ctx)
			if err != nil {
				errs <- err
				fmt.Println(err)
				return
			}
			msg := NewMsg(received.GetFrom(), received.GetData())
			fmt.Println(msg)
			err = b.append(msg)
			if err != nil {
				fmt.Println(err)
				errs <- err
				return
			}
		}
	}()

	// TODO: not pretty...
	select {
	case errS := <-errs:
		err = errS
	}
	if err != nil {
		return err
	}

	return nil
}

func (b *Box) Close() error {
	b.closeCh <- true
	close(b.closeCh)
	return nil
}

func (b *Box) GetMsgs() map[time.Time]*Msg {
	return b.msgs
}

func (b *Box) append(msg *Msg) error {
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
