package msg

import "github.com/h0n9/petit-chat/code"

type MsgType uint32

const (
	MsgTypeNone MsgType = iota + 1 // Msg 0 means something wrong
	MsgTypeRaw
	MsgTypeHelloSyn
	MsgTypeHelloAck
	MsgTypeBye // End of Subscription
	MsgTypeUpdateBox
)

var msgTypeMap = map[MsgType]string{
	MsgTypeNone:      "MsgTypeNone",
	MsgTypeRaw:       "MsgTypeRaw",
	MsgTypeHelloSyn:  "MsgTypeHelloSyn",
	MsgTypeHelloAck:  "MsgTypeHelloAck",
	MsgTypeBye:       "MsgTypeBye",
	MsgTypeUpdateBox: "MsgTypeUpdatgeBox",
}

func (mt MsgType) String() string {
	err := mt.Check()
	if err != nil {
		return ""
	}
	return msgTypeMap[mt]
}

func (mt MsgType) Check() error {
	_, exist := msgTypeMap[mt]
	if !exist {
		return code.UnknownMsgType
	}
	return nil
}
