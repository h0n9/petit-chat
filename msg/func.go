package msg

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

func (msg *Msg) execute(b *Box) error {
	mt := msg.GetType()
	mf, exist := msgFuncMap[mt]
	if !exist {
		mf = func(b *Box, m *Msg) error { return nil }
	}
	return mf(b, msg)
}
