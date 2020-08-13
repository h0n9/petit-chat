package msg

import (
	"time"

	"github.com/h0n9/petit-chat/code"
)

// Box refers to a chat room
type Box struct {
	msgs            map[time.Time]*Msg
	latestTimestamp time.Time
}

func NewBox() (*Box, error) {
	return &Box{msgs: make(map[time.Time]*Msg)}, nil
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
