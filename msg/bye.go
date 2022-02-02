package msg

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
)

type BodyBye struct {
	Persona *types.Persona `json:"persona"`
}

type Bye struct {
	Head
	Body BodyBye `json:"body"`
}

func NewMsgBye(peerID types.ID, clientAddr crypto.Addr, parentHash types.Hash, persona *types.Persona) *Msg {
	return NewMsg(&Bye{
		NewHead(peerID, clientAddr, parentHash, TypeBye),
		BodyBye{
			Persona: persona,
		},
	})
}

func (msg *Bye) GetBody() Body {
	return msg.Body
}

func (msg *Bye) Check(box *Box) error {
	auth := box.state.GetAuth()
	if !auth.IsPublic() && !auth.CanRead(msg.ClientAddr) {
		return code.NonReadPermission
	}
	if persona := box.getPersona(msg.ClientAddr); persona == nil {
		return code.NonExistingPersonaInBox
	}
	return nil
}

func (msg *Bye) Execute(box *Box) error {
	err := box.leave(msg.Body.Persona)
	if err != nil {
		return err
	}
	return nil
}
