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

func (body *BodyHelloSyn) Check(box *Box, addr crypto.Addr) error {
	if !box.auth.IsPublic() && !box.auth.CanRead(addr) {
		return code.NonReadPermission
	}
	return nil
}

func (body *BodyHelloSyn) Execute(box *Box, hash types.Hash) error {
	err := box.join(body.Persona)
	if err != nil {
		return err
	}

	// encrypt b.secretKey with msh.Persona.PubKey.GetKey()
	encryptedSecretKey, err := body.Persona.PubKey.Encrypt(box.secretKey.GetKey())
	if err != nil {
		return err
	}

	msg := NewMsg(box.myID, hash, &BodyHelloAck{
		Personae:           box.personae,
		Auth:               box.auth,
		EncryptedSecretKey: encryptedSecretKey,
	})

	err = box.Publish(msg, TypeHelloAck, false)
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

func (body *BodyHelloAck) Check(box *Box, addr crypto.Addr) error {
	return nil
}

func (body *BodyHelloAck) Execute(box *Box, hash types.Hash) error {
	secretKeyByte, err := box.myPrivKey.Decrypt(body.EncryptedSecretKey)
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
		box.personae = body.Personae
	}
	if util.HasField("auth", box) {
		box.auth = body.Auth
	}
	if util.HasField("secretKey", box) {
		box.secretKey = secretKey
	}

	return nil
}
