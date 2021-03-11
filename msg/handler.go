package msg

import (
	"github.com/h0n9/petit-chat/types"
)

type MsgHandler func(b *Box, psmsg *types.PubSubMsg) (bool, error)

func DefaultMsgHandler(b *Box, psmsg *types.PubSubMsg) (bool, error) {
	data := psmsg.GetData()
	msg, err := Decapsulate(data)
	if err != nil {
		return false, err
	}

	eos := msg.IsEOS() && (msg.GetFrom() == b.myID)

	// check if msg is proper and can be supported on protocol
	// improper msgs are dropped here
	err = msg.check()
	if err != nil {
		return eos, err
	}

	// execute msg with msgFunc
	msgFunc := msg.getMsgFunc()
	err = msgFunc(b, msg)
	if err != nil {
		return eos, err
	}

	// append msg
	readUntilIndex, err := b.append(msg)
	if err != nil {
		return eos, err
	}
	if psmsg.GetFrom() == b.myID {
		b.readUntilIndex = readUntilIndex
	} else {
		if b.msgSubCh != nil {
			b.msgSubCh <- msg
			b.readUntilIndex = readUntilIndex
		}
	}

	return eos, nil
}
