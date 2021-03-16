package msg

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/types"
)

type msgFunc func(b *Box, m *Msg) error

var msgFuncMap map[MsgType]msgFunc = map[MsgType]msgFunc{
	MsgTypeEOS: func(b *Box, m *Msg) error {
		if m.GetFrom() == b.myID {
			return nil
		}
		// TODO: remove msg.GetFrom() from b.Members
		return nil
	},
}

func (msg *Msg) check(b *Box) error {
	// check msgType
	mt := msg.GetType()
	err := mt.check()
	if err != nil {
		return err
	}

	// check msg.ParentMsgHash
	pm, err := msg.getParentMsg(b)
	if err != nil {
		return err
	}
	if pm != nil && types.IsHash(pm.ParentMsgHash) {
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
