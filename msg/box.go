package msg

import (
	"context"
	"fmt"
	"time"

	"github.com/h0n9/petit-chat/code"
)

// Box refers to a chat room
type Box struct {
	ctx        context.Context
	sub        *Sub
	subRoutine bool

	host  Peer
	peers []Peer
	msgs  map[time.Time]*Msg

	latestTimestamp time.Time
}

func NewBox(ctx context.Context, sub *Sub, host Peer, peers ...Peer) (
	*Box, error) {
	if sub == nil {
		return nil, code.ImproperSub
	}

	return &Box{
		sub:        sub,
		subRoutine: false,

		host:            host,
		peers:           peers,
		msgs:            make(map[time.Time]*Msg),
		latestTimestamp: time.Now(),
	}, nil
}

func (b *Box) GetHost() Peer {
	return b.host
}

func (b *Box) GetPeers() []Peer {
	return b.peers
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

func (b *Box) turnOnSubRoutine() {
	b.subRoutine = true

	go func() {
		for {
			rawMsg, err := b.sub.Next(b.ctx)
			if err != nil {
				continue
			}
			if len(rawMsg.Data) == 0 {
				continue
			}
			if rawMsg.GetFrom().String() == b.host.GetID().String() {
				continue
			}

			msg, err := UnmarshalJSON(rawMsg.GetData())
			if err != nil {
				fmt.Printf("%s: %s\n", err, rawMsg.GetData())
				continue
			}

			fmt.Printf("\x1b[32m%s: %s\x1b[0m\n> ",
				msg.GetFrom().id,
				msg.GetData(),
			)
		}
	}()
}
