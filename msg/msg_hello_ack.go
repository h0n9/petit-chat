package msg

import (
	"encoding/json"

	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type MsgStructHelloAck struct {
	Persona            *types.Persona `json:"persona"`
	Auth               *types.Auth    `json:"auth"`
	EncryptedSecretKey []byte         `json:"encrypted_secret_key"`
}

func NewMsgStructHelloAck(
	persona *types.Persona, auth *types.Auth, encryptedSecretKey []byte,
) *MsgStructHelloAck {
	return &MsgStructHelloAck{
		Persona:            persona,
		Auth:               auth,
		EncryptedSecretKey: encryptedSecretKey,
	}
}

func (msha *MsgStructHelloAck) Encapsulate() ([]byte, error) {
	return json.Marshal(msha)
}

func (msha *MsgStructHelloAck) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msha)
}

func (msha *MsgStructHelloAck) Execute(b *Box, fromPeerID types.ID) error {
	err := b.join(msha.Persona)
	if err != nil {
		return err
	}

	if fromPeerID == b.myID {
		return nil
	}

	// back msg
	// decrypt msh.encrypted
	secretKeyByte, err := b.myPrivKey.Decrypt(msha.EncryptedSecretKey)
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
		b.auth = msha.Auth
	}

	return nil
}
