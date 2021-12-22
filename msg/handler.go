package msg

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/types"
)

type MsgHandler func(b *Box, msg *Msg) (bool, error)

func DefaultMsgHandler(b *Box, msg *Msg) (bool, error) {
	eos := msg.IsEOS() && (msg.GetPeerID() == b.myID)

	// msg handling flow:
	//   check -> append -> execute -> (received)

	// check if msg is proper and can be supported on protocol
	// improper msgs are dropped here
	err := msg.check(b)
	if err != nil {
		return eos, err
	}

	hash := msg.GetHash()
	addr := msg.Signature.PubKey.Address()

	// check msg.Body
	err = msg.Body.Check(b, hash, addr)
	if err != nil && err != code.SelfMsg {
		return eos, err
	}

	// append msg
	readUntilIndex, err := b.append(msg)
	if err != nil {
		return eos, err
	}

	// execute msg.Body
	err = msg.Body.Execute(b, hash, addr)
	if err != nil {
		return eos, err
	}

	canRead := b.msgSubCh != nil
	if canRead {
		b.msgSubCh <- msg
		b.readUntilIndex = readUntilIndex
	}
	if msg.Type != TypeMeta {
		msgMeta := NewMsg(b.myID, msg.Hash, TypeMeta, &BodyMeta{
			Meta: types.NewMeta(true, canRead, false),
		})
		err := b.Publish(msgMeta, true)
		if err != nil {
			return eos, err
		}
	}

	return eos, nil
}
