package msg

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/types"
)

// TODO: fix size constraint
const (
	MaxDataSize     = 5000
	MaxMetadataSize = 1000
)

type BodyRaw struct {
	Data     []byte `json:"data"`
	Metadata []byte `json:"metadata"`
}

type Raw struct {
	Head
	Body BodyRaw `json:"body"`
}

func NewMsgRaw(box *Box, parentHash types.Hash, data []byte, metadata []byte) *Msg {
	return NewMsg(&Raw{
		Head: NewHead(box, parentHash, TypeRaw),
		Body: BodyRaw{
			Data:     data,
			Metadata: metadata,
		},
	})
}

func (msg *Raw) GetBody() Body {
	return msg.Body
}

func (msg *Raw) Check(box *Box) error {
	clientAddr := msg.GetClientAddr()
	if !box.auth.CanWrite(clientAddr) {
		return code.NonWritePermission
	}
	if len(msg.Body.Data) > MaxDataSize {
		return code.TooBigMsgData
	}
	if len(msg.Body.Metadata) > MaxMetadataSize {
		return code.TooBigMsgMetadata
	}
	if !box.auth.CanWrite(clientAddr) {
		return code.NonWritePermission
	}
	return nil
}

func (msg *Raw) Execute(box *Box) error {
	return nil
}
