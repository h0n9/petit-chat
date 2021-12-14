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

func (body *BodyRaw) Check(box *Box, from *From) error {
	if len(body.Data) > MaxDataSize {
		return code.TooBigMsgData
	}
	if len(body.Metadata) > MaxMetadataSize {
		return code.TooBigMsgMetadata
	}
	if !box.auth.CanWrite(from.ClientAddr) {
		return code.NonWritePermission
	}
	return nil
}

func (body *BodyRaw) Execute(box *Box, hash types.Hash) error {
	return nil
}
