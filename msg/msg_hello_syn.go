package msg

import (
	"encoding/json"

	"github.com/h0n9/petit-chat/types"
)

type MsgStructHelloSyn struct {
	Persona *types.Persona `json:"persona"`
}

func NewMsgStructHelloSyn(persona *types.Persona) *MsgStructHelloSyn {
	return &MsgStructHelloSyn{Persona: persona}
}

func (mshs *MsgStructHelloSyn) Encapsulate() ([]byte, error) {
	return json.Marshal(mshs)
}

func (mshs *MsgStructHelloSyn) Decapsulate(data []byte) error {
	return json.Unmarshal(data, mshs)
}

func (mshs *MsgStructHelloSyn) Execute(b *Box, fromPeerID types.ID, hash types.Hash) error {
	// TODO: check more constraints for fromPeerID
	if fromPeerID == b.myID {
		return nil
	}

	err := b.join(mshs.Persona)
	if err != nil {
		return err
	}

	// encrypt b.secretKey with msh.Persona.PubKey.GetKey()
	encryptedSecretKey, err := mshs.Persona.PubKey.Encrypt(b.secretKey.GetKey())
	if err != nil {
		return err
	}

	msh := NewMsgStructHelloAck(b.myPersona, b.auth, encryptedSecretKey)
	data, err := msh.Encapsulate()
	if err != nil {
		return err
	}

	err = b.Publish(MsgTypeHelloAck, hash, false, data)
	if err != nil {
		return err
	}

	return nil
}
