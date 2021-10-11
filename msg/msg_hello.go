package msg

import (
	"encoding/json"

	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type MsgStructHello struct {
	Persona            *types.Persona `json:"persona"`
	Auth               *types.Auth    `json:"auth"`
	EncryptedSecretKey []byte         `json:"encrypted_secret_key"`
}

func NewMsgStructHello(persona *types.Persona, auth *types.Auth, encryptedSecretKey []byte) *MsgStructHello {
	return &MsgStructHello{
		Persona:            persona,
		Auth:               auth,
		EncryptedSecretKey: encryptedSecretKey,
	}
}

func (msh *MsgStructHello) Encapsulate() ([]byte, error) {
	return json.Marshal(msh)
}

func (msh *MsgStructHello) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msh)
}

func (msh *MsgStructHello) Execute(b *Box, m *Msg) error {
	err := b.join(msh.Persona)
	if err != nil {
		return err
	}

	if m.GetFrom().PeerID == b.myID {
		return nil
	}

	if m.ParentMsgHash.IsEmpty() {
		// new msg
		pmhash, err := m.Hash()
		if err != nil {
			return err
		}

		// encrypt b.secretKey with msh.Persona.PubKey.GetKey()
		encryptedSecretKey, err := msh.Persona.PubKey.Encrypt(b.secretKey.GetKey())
		if err != nil {
			return err
		}

		msh := NewMsgStructHello(b.myPersona, b.auth, encryptedSecretKey)
		data, err := msh.Encapsulate()
		if err != nil {
			return err
		}

		err = b.Publish(MsgTypeHello, pmhash, false, data)
		if err != nil {
			return err
		}

		return nil
	}

	// back msg
	// decrypt msh.encrypted
	secretKeyByte, err := b.myPrivKey.Decrypt(msh.EncryptedSecretKey)
	if err != nil {
		// TODO: handle or log error somehow
		// this could not be a real error
		return nil
	}
	secretKey, err := crypto.NewSecretKey(secretKeyByte)
	if err != nil {
		return err
	}

	// apply to msgBox struct values
	if util.HasField("secretKey", b) {
		b.secretKey = secretKey
	}
	if util.HasField("auth", b) {
		b.auth = msh.Auth
	}

	return nil
}
