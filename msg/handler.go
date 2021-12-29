package msg

import (
	"github.com/h0n9/petit-chat/code"
)

type MsgHandler func(box *Box, msg *Msg) (bool, error)

func DefaultMsgHandler(box *Box, msg *Msg) (bool, error) {
	eos := msg.IsEOS() && (msg.GetPeerID() == box.myID)

	// msg handling flow:
	//   check -> execute -> append

	// check if msg is proper and can be supported on protocol
	// improper msgs are dropped here
	err := msg.check(box)
	if err != nil {
		return eos, err
	}

	// check msg.Body
	err = msg.Check(box)
	if err != nil && err != code.SelfMsg {
		return eos, err
	}

	// execute msg.Body
	err = msg.Execute(box)
	if err != nil {
		return eos, err
	}

	// append msg
	readUntilIndex, err := box.append(msg)
	if err != nil {
		return eos, err
	}

	if msg.GetPeerID() == box.myID {
		box.readUntilIndex = readUntilIndex
	} else {
		if box.msgSubCh != nil {
			box.msgSubCh <- msg
			box.readUntilIndex = readUntilIndex
		}
	}

	return eos, nil
}
