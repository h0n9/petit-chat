package msg

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/types"
)

type msgFunc func(b *Box, m *Msg) error

var msgFuncMap map[types.Msg]msgFunc = map[types.Msg]msgFunc{
	types.MsgEOS: func(b *Box, m *Msg) error {
		if m.GetFrom() == b.myID {
			return nil
		}
		// TODO: remove msg.GetFrom() from b.Members
		return nil
	},
	types.MsgNewbie: func(b *Box, m *Msg) error {
		return nil
	},
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
