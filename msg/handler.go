package msg

import (
	"github.com/h0n9/petit-chat/types"
)

type MsgHandler func(b *Box, psmsg *types.PubSubMsg) error

func DefaultMsgHandler(b *Box, pbmsg *types.PubSubMsg) error {
	data := pbmsg.GetData()
	msg, err := Decapsulate(data)
	if err != nil {
		return err
	}
	// TODO: consider if this a right way to handle closing subscription
	if msg.IsEOS() {
		if pbmsg.GetFrom() == b.myID {
			b.sub.Cancel()
			err := b.topic.Close()
			if err != nil {
				return err
			}
			b.sub = nil
		}
		return nil
	}
	readUntilIndex, err := b.append(msg)
	if err != nil {
		return err
	}
	if pbmsg.GetFrom() == b.myID {
		b.readUntilIndex = readUntilIndex
	} else {
		if b.msgSubCh != nil {
			b.msgSubCh <- msg
			b.readUntilIndex = readUntilIndex
		}
	}
	return nil
}
