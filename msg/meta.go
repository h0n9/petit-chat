package msg

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/types"
)

type BodyMeta struct {
	TargetMsgHash types.Hash `json:"target_msg_hash"`
	Meta          types.Meta `json:"meta"`
}

type Meta struct {
	Head
	Body BodyMeta `json:"body"`
}

func NewMsgMeta(box *Box, parentHash types.Hash, targetMsgHash types.Hash, meta types.Meta) *Msg {
	return NewMsg(&Meta{
		Head: NewHead(box, parentHash, TypeMeta),
		Body: BodyMeta{
			TargetMsgHash: targetMsgHash,
			Meta:          meta,
		},
	})
}

func (msg *Meta) GetBody() Body {
	return msg.Body
}

func (msg *Meta) Check(box *Box) error {
	clientAddr := msg.GetClientAddr()
	if msg.Body.Meta.Received() || msg.Body.Meta.Read() {
		if !box.auth.IsPublic() && !box.auth.CanRead(clientAddr) {
			return code.NonReadPermission
		}
		if msg.Body.TargetMsgHash.IsEmpty() {
			return code.UnknownMsgType
		}
	}
	if msg.Body.Meta.Typing() {
		if !box.auth.CanWrite(clientAddr) {
			return code.NonWritePermission
		}
		if !msg.Body.TargetMsgHash.IsEmpty() {
			return code.UnknownMsgType
		}
	}
	return nil
}

func (msg *Meta) Execute(box *Box) error {
	if msg.Body.Meta.Received() || msg.Body.Meta.Read() {
		targetMsg := box.GetMsg(msg.Body.TargetMsgHash)
		if targetMsg == nil {
			return code.NonExistingMsg
		}
		targetMsg.MergeMeta(msg.GetClientAddr(), msg.Body.Meta)
	}
	// if msg.Body.Meta.Typing() {
	// 	// TODO: do something
	// }
	return nil
}
