package msg

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
)

type msgFunc func(b *Box, m *Msg) error

var msgFuncMap map[MsgType]msgFunc = map[MsgType]msgFunc{
	MsgTypeHello: msgFuncHello,
	MsgTypeBye:   msgFuncBye,
}

func (msg *Msg) check(b *Box) error {
	// check msgType
	mt := msg.GetType()
	err := mt.Check()
	if err != nil {
		return err
	}

	// check msg.ParentMsgHash
	pm, err := msg.getParentMsg(b)
	if err != nil {
		return err
	}
	if pm != nil && !types.IsEmpty(pm.ParentMsgHash) {
		return code.AlreadyHavingParentMsg
	}

	// TODO: add more constraints

	return nil
}

func (msg *Msg) execute(b *Box) error {
	mt := msg.GetType()
	mf, exist := msgFuncMap[mt]
	if !exist {
		mf = func(b *Box, m *Msg) error { return nil }
	}
	return mf(b, msg)
}

func msgFuncHello(b *Box, m *Msg) error {
	msh := NewMsgStructHello(nil, nil, nil)
	err := msh.Decapsulate(m.GetData())
	if err != nil {
		return err
	}
	err = b.join(msh.Persona)
	if err != nil {
		return err
	}

	if m.GetFrom().PeerID == b.myID {
		return nil
	}

	if types.IsEmpty(m.ParentMsgHash) {
		// new msg
		pmhash, err := m.Hash()
		if err != nil {
			return err
		}

		// encrypt b.secretKey with msh.Persona.PubKey.GetKey()
		encryptedSecretKey, err := msh.Persona.PubKey.Encrypt(b.secretKey.GetKey())
		if err != nil {
			return err
		}

		msh := NewMsgStructHello(b.myPersona, b.auth, encryptedSecretKey)
		data, err := msh.Encapsulate()
		if err != nil {
			return err
		}

		err = b.Publish(MsgTypeHello, pmhash, false, data)
		if err != nil {
			return err
		}

		return nil
	}

	// back msg
	// decrypt msh.encrypted
	secretKeyByte, err := b.myPrivKey.Decrypt(msh.EncryptedSecretKey)
	if err != nil {
		return err
	}
	secretKey, err := crypto.NewSecretKey(secretKeyByte)
	if err != nil {
		return err
	}
	b.secretKey = secretKey

	return nil
}

func msgFuncBye(b *Box, m *Msg) error {
	if m.GetFrom().PeerID == b.myID {
		return nil
	}

	msb := NewMsgStructBye(nil)
	err := msb.Decapsulate(m.GetData())
	if err != nil {
		return err
	}

	err = b.leave(msb.Persona)
	if err != nil {
		return err
	}

	return nil
}
