package msg

import (
	"github.com/h0n9/petit-chat/code"
)

type MsgType uint32

const (
	MsgTypeNone MsgType = iota + 1 // MsgType 0 means something wrong
	MsgTypeText
	MsgTypeImage
	MsgTypeVideo
	MsgTypeAudio
	MsgTypeRaw
	MsgTypeEOS // End of Subscription
)

var msgTypeMap = map[MsgType]string{
	MsgTypeNone:  "MsgTypeNone",
	MsgTypeText:  "MsgTypeText",
	MsgTypeImage: "MsgTypeImage",
	MsgTypeVideo: "MsgTypeVideo",
	MsgTypeAudio: "MsgTypeAudio",
	MsgTypeRaw:   "MsgTypeRaw",
	MsgTypeEOS:   "MsgTypeEOS",
}

func (mt MsgType) String() string {
	err := mt.check()
	if err != nil {
		return ""
	}
	return msgTypeMap[mt]
}

func (mt MsgType) Check() error {
	return mt.check()
}

func (mt MsgType) check() error {
	_, exist := msgTypeMap[mt]
	if !exist {
		return code.UnknownMsgType
	}
	return nil
}
