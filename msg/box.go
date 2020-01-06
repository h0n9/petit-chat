package msg

import (
	"time"

	"github.com/h0n9/petit-chat/code"
)

type MsgBox struct {
	sub Sub

	Nickname string
	Peers    []Peer
	Msgs     map[time.Time]*Msg

	LatestTimestamp time.Time
}

func NewMsgBox(nickname string, peers []Peer) *MsgBox {
	return &MsgBox{
		Nickname:        nickname,
		Peers:           peers,
		Msgs:            map[time.Time]*Msg{},
		LatestTimestamp: time.Now(),
	}
}

func (mb *MsgBox) Append(timestamp time.Time, msg *Msg) error {
	_, exist := mb.Msgs[timestamp]
	if exist {
		return code.AlreadyAppendedMsg
	}

	if mb.LatestTimestamp.Before(timestamp) {
		mb.LatestTimestamp = timestamp
	}

	mb.Msgs[timestamp] = msg

	return nil
}
