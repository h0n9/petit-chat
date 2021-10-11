package msg

import (
	"encoding/json"

	"github.com/h0n9/petit-chat/types"
)

type MsgStructBye struct {
	Persona *types.Persona `json:"persona"`
}

func NewMsgStructBye(persona *types.Persona) *MsgStructBye {
	return &MsgStructBye{Persona: persona}
}

func (msb *MsgStructBye) Encapsulate() ([]byte, error) {
	return json.Marshal(msb)
}

func (msb *MsgStructBye) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msb)
}

func (msb *MsgStructBye) Execute(b *Box, m *Msg) error {
	if m.GetFrom().PeerID == b.myID {
		return nil
	}

	err := b.leave(msb.Persona)
	if err != nil {
		return err
	}

	return nil
}
