package msg

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/types"
)

type BodyBye struct {
	Persona *types.Persona `json:"persona"`
}

func (body *BodyBye) Check(box *Box, from *From) error {
	// if from.PeerID == box.myID {
	// 	return code.SelfMsg
	// }
	if persona := box.getPersona(from.ClientAddr); persona == nil {
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
