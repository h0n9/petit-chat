package msg

import (
	"github.com/h0n9/petit-chat/code"
)

type MsgHandler func(b *Box, msg *Msg) (bool, error)

func DefaultMsgHandler(b *Box, msg *Msg) (bool, error) {
	eos := msg.IsEOS() && (msg.GetPeerID() == b.myID)

	// msg handling flow:
	//   check -> execute -> append -> (received)

	// check if msg is proper and can be supported on protocol
	// improper msgs are dropped here
	err := msg.check(b)
	if err != nil {
		return eos, err
	}

	addr := msg.Signature.PubKey.Address()
	hash := msg.GetHash()

	// check msg.Body
	err = msg.Body.Check(b, addr)
	if err != nil && err != code.SelfMsg {
		return eos, err
	}

	// execute msg.Body
	err = msg.Body.Execute(b, hash)
	if err != nil {
		return eos, err
	}

	// append msg
	readUntilIndex, err := b.append(msg)
	if err != nil {
		return eos, err
	}

	if msg.GetPeerID() == b.myID {
		b.readUntilIndex = readUntilIndex
	} else {
		if b.msgSubCh != nil {
			b.msgSubCh <- msg
			b.readUntilIndex = readUntilIndex
		}
	}

	return eos, nil
}
