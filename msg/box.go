package msg

import (
	"context"
	"fmt"
	"time"

	"github.com/h0n9/petit-chat/code"
)

// MsgBox refers to a chat room
type MsgBox struct {
	ctx        context.Context
	sub        *Sub
	subRoutine bool

	host  Peer
	peers []Peer
	msgs  map[time.Time]*Msg

	latestTimestamp time.Time
}

func NewMsgBox(ctx context.Context, sub *Sub, host Peer, peers ...Peer) (
	*MsgBox, error) {
	if sub == nil {
		return nil, code.ImproperSub
	}

	return &MsgBox{
		sub:        sub,
		subRoutine: false,

		host:            host,
		peers:           peers,
		msgs:            make(map[time.Time]*Msg),
		latestTimestamp: time.Now(),
	}, nil
}

func (mb *MsgBox) GetHost() Peer {
	return mb.host
}

func (mb *MsgBox) GetPeers() []Peer {
	return mb.peers
}

func (mb *MsgBox) Append(msg *Msg) error {
	_, exist := mb.msgs[msg.Timestamp]
	if exist {
		return code.AlreadyAppendedMsg
	}

	if mb.latestTimestamp.Before(msg.Timestamp) {
		mb.latestTimestamp = msg.Timestamp
	}

	mb.msgs[msg.Timestamp] = msg

	return nil
}

func (mb *MsgBox) turnOnSubRoutine() {
	mb.subRoutine = true

	go func() {
		for {
			rawMsg, err := mb.sub.Next(mb.ctx)
			if err != nil {
				continue
			}
			if len(rawMsg.Data) == 0 {
				continue
			}
			if rawMsg.GetFrom().String() == mb.host.ID.String() {
				continue
			}

			msg, err := UnmarshalJSON(rawMsg.GetData())
			if err != nil {
				fmt.Printf("%s: %s\n", err, rawMsg.GetData())
				continue
			}

			fmt.Printf("\x1b[32m%s: %s\x1b[0m\n> ",
				msg.GetFrom(),
				msg.GetData(),
			)
		}
	}()
}
