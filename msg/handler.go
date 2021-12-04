package msg

import (
	"github.com/h0n9/petit-chat/code"
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
	case TypeRaw:
		var msr MsgStructRaw
		err := msr.Decapsulate(msg.Data)
		if err != nil {
			return eos, err
		}
		if !b.auth.CanWrite(from.ClientAddr) {
			return eos, code.NonWritePermission
		}
	case TypeHelloSyn:
		var mshs MsgStructHelloSyn
		err := mshs.Decapsulate(msg.Data)
		if err != nil {
			return eos, err
		}
		// private and cannot read
		if !b.auth.IsPublic() && !b.auth.CanRead(from.ClientAddr) {
			return eos, code.NonReadPermission
		}
		err = mshs.Execute(b, from.PeerID, hash)
		if err != nil {
			return eos, err
		}
	case TypeHelloAck:
		var msha MsgStructHelloAck
		err := msha.Decapsulate(msg.Data)
		if err != nil {
			return eos, err
		}
		// ok, err := b.auth.CheckMinPerm(from.ClientAddr)
		// if err != nil {
		// 	return eos, err
		// }
		// if !ok {
		// 	return eos, code.NonMinimumPermission
		// }
		err = msha.Execute(b, from.PeerID)
		if err != nil {
			return eos, err
		}
	case TypeBye:
		var msb MsgStructBye
		err := msb.Decapsulate(msg.Data)
		if err != nil {
			return eos, err
		}
		// private and cannot read
		// if !b.auth.IsPublic() && !b.auth.CanRead(from.ClientAddr) {
		// 	return eos, code.NonReadPermission
		// }
		err = msb.Execute(b, from.PeerID)
		if err != nil {
			return eos, err
		}
	case TypeUpdateSyn:
		var msus MsgStructUpdateSyn
		err := msus.Decapsulate(msg.Data)
		if err != nil {
			return eos, err
		}
		if b.auth.CanExecute(from.ClientAddr) {
			err = msus.Execute(b)
			if err != nil {
				return eos, err
			}
		}
		msua := NewMsgStructUpdateAck(b.auth, b.personae)
		data, err := msua.Encapsulate()
		if err != nil {
			return eos, err
		}
		err = b.Publish(TypeUpdateAck, types.Hash{}, true, data)
		if err != nil {
			return eos, err
		}
	case TypeUpdateAck:
		var msua MsgStructUpdateAck
		err := msua.Decapsulate(msg.Data)
		if err != nil {
			return eos, err
		}
		if b.auth.CanExecute(from.ClientAddr) {
			err = msua.Execute(b)
			if err != nil {
				return eos, err
			}
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
