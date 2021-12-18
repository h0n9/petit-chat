package msg

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
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

func (body *BodyRaw) Check(box *Box, addr crypto.Addr) error {
	if !box.auth.CanWrite(addr) {
		return code.NonWritePermission
	}
	if len(body.Data) > MaxDataSize {
		return code.TooBigMsgData
	}
	if len(body.Metadata) > MaxMetadataSize {
		return code.TooBigMsgMetadata
	}
	if !box.auth.CanWrite(addr) {
		return code.NonWritePermission
	}
	return nil
}

func (body *BodyRaw) Execute(box *Box, hash types.Hash) error {
	return nil
}
