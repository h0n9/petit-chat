package msg

import (
	"github.com/h0n9/petit-chat/types"
)

type MsgHandler func(b *Box, msg *Msg, fromID types.ID) (bool, error)

func DefaultMsgHandler(b *Box, msg *Msg, fromID types.ID) (bool, error) {
	eos := msg.IsEOS() && (msg.GetFrom().PeerID == b.myID)

	// msg handling flow:
	//   check -> decrypt(optional) -> decapsulate and execute(optional) -> append

	// check if msg is proper and can be supported on protocol
	// improper msgs are dropped here
	err := msg.check(b)
	if err != nil {
		return eos, err
	}

	// decrypt if encrypted
	if msg.Encrypted {
		decryptedData, err := b.secretKey.Decrypt(msg.GetData())
		if err != nil {
			return eos, err
		}
		msg.SetData(decryptedData)
	}

	from := msg.GetFrom()
	hash := msg.GetHash()

	// decapsulate and execute
	switch msg.Type {
	case MsgTypeHelloSyn:
		mshs := NewMsgStructHelloSyn(nil)
		err := mshs.Decapsulate(msg.Data)
		if err != nil {
			return eos, err
		}
		err = mshs.Execute(b, from.PeerID, hash)
		if err != nil {
			return eos, err
		}
	case MsgTypeHelloAck:
		msha := NewMsgStructHelloAck(nil, nil, nil)
		err := msha.Decapsulate(msg.Data)
		if err != nil {
			return eos, err
		}
		err = msha.Execute(b, from.PeerID)
		if err != nil {
			return eos, err
		}
	case MsgTypeBye:
		msb := NewMsgStructBye(nil)
		err := msb.Decapsulate(msg.Data)
		if err != nil {
			return eos, err
		}
		err = msb.Execute(b, from.PeerID)
		if err != nil {
			return eos, err
		}
	case MsgTypeUpdateBox:
		msub := NewMsgStructUpdateBox(nil, nil)
		err := msub.Decapsulate(msg.Data)
		if err != nil {
			return eos, err
		}
		err = msub.Execute(b)
		if err != nil {
			return eos, err
		}
	}

	// append msg
	readUntilIndex, err := b.append(msg)
	if err != nil {
		return eos, err
	}
	if fromID == b.myID {
		b.readUntilIndex = readUntilIndex
	} else {
		if b.msgSubCh != nil {
			b.msgSubCh <- msg
			b.readUntilIndex = readUntilIndex
		}
	}

	return eos, nil
}
