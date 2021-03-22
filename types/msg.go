package types

import (
	"github.com/h0n9/petit-chat/code"
)

type Msg uint32

const (
	MsgNone Msg = iota + 1 // Msg 0 means something wrong
	MsgText
	MsgImage
	MsgVideo
	MsgAudio
	MsgRaw
	MsgEOS // End of Subscription
	MsgNewbie
)

var msgMap = map[Msg]string{
	MsgNone:  "MsgNone",
	MsgText:  "MsgText",
	MsgImage: "MsgImage",
	MsgVideo: "MsgVideo",
	MsgAudio: "MsgAudio",
	MsgRaw:   "MsgRaw",
	MsgEOS:   "MsgEOS",
	MsgNewbie: "MsgNewbie",
}

func (mt Msg) String() string {
	err := mt.Check()
	if err != nil {
		return ""
	}
	return msgMap[mt]
}

func (mt Msg) Check() error {
	_, exist := msgMap[mt]
	if !exist {
		return code.UnknownMsgType
	}
	return nil
}
