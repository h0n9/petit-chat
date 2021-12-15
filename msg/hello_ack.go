package msg

import (
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

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
