package msg

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
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

func NewMsgMeta(peerID types.ID, clientAddr crypto.Addr, parentHash, targetMsgHash types.Hash, meta types.Meta) *Msg {
	return NewMsg(&Meta{
		NewHead(peerID, clientAddr, parentHash, TypeMeta),
		BodyMeta{
			TargetMsgHash: targetMsgHash,
			Meta:          meta,
		},
	})
}

func (msg *Meta) GetBody() Body {
	return msg.Body
}

func (msg *Meta) Check(hash types.Hash, helper Helper) error {
	state := helper.GetState()
	auth := state.GetAuth()

	clientAddr := msg.GetClientAddr()
	if msg.Body.Meta.Received() || msg.Body.Meta.Read() {
		if !auth.IsPublic() && !auth.CanRead(clientAddr) {
			return code.NonReadPermission
		}
		if msg.Body.TargetMsgHash.IsEmpty() {
			return code.UnknownMsgType
		}
	}
	if msg.Body.Meta.Typing() {
		if !auth.CanWrite(clientAddr) {
			return code.NonWritePermission
		}
		if !msg.Body.TargetMsgHash.IsEmpty() {
			return code.UnknownMsgType
		}
	}
	return nil
}

func (msg *Meta) Execute(hash types.Hash, helper Helper) error {
	state := helper.GetState()
	targetMsgHash := msg.Body.TargetMsgHash
	clientAddr := msg.GetClientAddr()
	meta := msg.Body.Meta
	if meta.Received() || meta.Read() {
		state.UpdateMeta(targetMsgHash, clientAddr, meta)
	}
	if msg.Body.Meta.Typing() {
		// TODO: do something
	}
	return nil
}
