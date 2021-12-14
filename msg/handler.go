package msg

import "github.com/h0n9/petit-chat/code"

type MsgHandler func(b *Box, msg *Msg) (bool, error)

func DefaultMsgHandler(b *Box, msg *Msg) (bool, error) {
	eos := msg.IsEOS() && (msg.GetFrom().PeerID == b.myID)

	// msg handling flow:
	//   check -> execute -> append

	// check if msg is proper and can be supported on protocol
	// improper msgs are dropped here
	err := msg.check(b)
	if err != nil {
		return eos, err
	}

	from := msg.GetFrom()
	hash := msg.GetHash()

	// check msg.Body
	err = msg.Body.Check(b, from)
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

	if from.PeerID == b.myID {
		b.readUntilIndex = readUntilIndex
	} else {
		if b.msgSubCh != nil {
			b.msgSubCh <- msg
			b.readUntilIndex = readUntilIndex
		}
	}

	return eos, nil
}
