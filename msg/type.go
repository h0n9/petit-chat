package msg

import "github.com/h0n9/petit-chat/code"

type MsgType uint32

const (
	MsgTypeNone MsgType = iota + 1 // Msg 0 means something wrong
	MsgTypeText
	MsgTypeImage
	MsgTypeVideo
	MsgTypeAudio
	MsgTypeRaw
	MsgTypeHello
	MsgTypeBye // End of Subscription
)

var msgTypeMap = map[MsgType]string{
	MsgTypeNone:  "MsgTypeNone",
	MsgTypeText:  "MsgTypeText",
	MsgTypeImage: "MsgTypeImage",
	MsgTypeVideo: "MsgTypeVideo",
	MsgTypeAudio: "MsgTypeAudio",
	MsgTypeRaw:   "MsgTypeRaw",
	MsgTypeHello: "MsgTypeHello",
	MsgTypeBye:   "MsgTypeBye",
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
