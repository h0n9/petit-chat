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

	Nickname string
	Host     Peer
	Peers    []Peer
	Msgs     map[time.Time]*Msg

	LatestTimestamp time.Time
}

func NewMsgBox(ctx context.Context, sub *Sub, host Peer,
	nickname string, peers ...Peer) (*MsgBox, error) {
	if sub == nil {
		return nil, code.ImproperSub
	}

	return &MsgBox{
		sub:        sub,
		subRoutine: false,

		Nickname:        nickname,
		Host:            host,
		Peers:           peers,
		Msgs:            make(map[time.Time]*Msg),
		LatestTimestamp: time.Now(),
	}, nil
}

func (mb *MsgBox) GetPeers() []Peer {
	return mb.Peers
}

func (mb *MsgBox) Append(msg *Msg) error {
	_, exist := mb.Msgs[msg.Timestamp]
	if exist {
		return code.AlreadyAppendedMsg
	}

	if mb.LatestTimestamp.Before(msg.Timestamp) {
		mb.LatestTimestamp = msg.Timestamp
	}

	mb.Msgs[msg.Timestamp] = msg

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
			if rawMsg.GetFrom().String() == mb.Host.ID.String() {
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
