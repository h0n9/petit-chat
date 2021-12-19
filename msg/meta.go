package msg

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
)

type BodyMeta struct {
	Meta types.Meta `json:"meta"`
}

func (body *BodyMeta) Check(box *Box, addr crypto.Addr) error {
	if body.Meta.Received() || body.Meta.Read() {
		if !box.auth.IsPublic() && !box.auth.CanRead(addr) {
			return code.NonReadPermission
		}
	}
	if body.Meta.Typing() && !box.auth.CanWrite(addr) {
		return code.NonWritePermission
	}
	return nil
}

func (body *BodyMeta) Execute(box *Box, hash types.Hash) error {
	return nil
}
