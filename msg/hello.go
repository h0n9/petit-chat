package msg

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type BodyHelloSyn struct {
	Persona *types.Persona `json:"persona"`
}

type HelloSyn struct {
	Head
	Body BodyHelloSyn `json:"body"`
}

func NewMsgHelloSyn(box *Box, parentHash types.Hash, persona *types.Persona) *Msg {
	return NewMsg(&HelloSyn{
		NewHead(box, parentHash, TypeHelloSyn),
		BodyHelloSyn{
			Persona: persona,
		},
	})
}

func (msg *HelloSyn) GetBody() Body {
	return msg.Body
}

func (msg *HelloSyn) Check(box *Box) error {
	if !box.auth.IsPublic() && !box.auth.CanRead(msg.GetClientAddr()) {
		return code.NonReadPermission
	}
	return nil
}

func (msg *HelloSyn) Execute(box *Box) error {
	err := box.join(msg.Body.Persona)
	if err != nil {
		return err
	}

	encryptedSecretKey, err := msg.Body.Persona.PubKey.Encrypt(box.secretKey.GetKey())
	if err != nil {
		return err
	}

	msgAck := NewMsgHelloAck(box, Hash(msg), box.personae, box.auth, encryptedSecretKey)
	err = box.Publish(msgAck, false)
	if err != nil {
		return err
	}

	return nil
}

type BodyHelloAck struct {
	Personae           types.Personae `json:"personae"`
	Auth               *types.Auth    `json:"auth"`
	EncryptedSecretKey []byte         `json:"encrypted_secret_key"`
}

type HelloAck struct {
	Head
	Body BodyHelloAck `json:"body"`
}

func NewMsgHelloAck(box *Box, parentHash types.Hash,
	personae types.Personae, auth *types.Auth, encryptedSecretKey []byte) *Msg {
	return NewMsg(&HelloAck{
		NewHead(box, parentHash, TypeHelloAck),
		BodyHelloAck{
			Personae:           personae,
			Auth:               auth,
			EncryptedSecretKey: encryptedSecretKey,
		},
	})
}

func (msg *HelloAck) GetBody() Body {
	return msg.Body
}

func (msg HelloAck) Check(box *Box) error {
	parentMsg, err := msg.getParentMsg(box)
	if err != nil {
		return err
	}
	if parentMsg == nil {
		return code.NonExistingParent
	}
	return nil
}

func (msg HelloAck) Execute(box *Box) error {
	secretKeyByte, err := box.myPrivKey.Decrypt(msg.Body.EncryptedSecretKey)
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
	if util.HasField("personae", box) {
		box.personae = msg.Body.Personae
	}
	if util.HasField("auth", box) {
		box.auth = msg.Body.Auth
	}
	if util.HasField("secretKey", box) {
		box.secretKey = secretKey
	}

	return nil
}
