package msg

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
)

type BodyBye struct {
	Persona *types.Persona `json:"persona"`
}

func (body *BodyBye) Check(box *Box, addr crypto.Addr) error {
	if !box.auth.IsPublic() && !box.auth.CanRead(addr) {
		return code.NonReadPermission
	}
	if persona := box.getPersona(addr); persona == nil {
		return code.NonExistingPersonaInBox
	}
	return nil
}

func (body *BodyBye) Execute(box *Box, hash types.Hash) error {
	err := box.leave(body.Persona)
	if err != nil {
		return err
	}
	return nil
}
