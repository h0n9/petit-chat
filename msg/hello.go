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

func NewMsgHelloSyn(peerID types.ID, clientAddr crypto.Addr, parentHash types.Hash, persona *types.Persona) *Msg {
	return NewMsg(&HelloSyn{
		NewHead(peerID, clientAddr, parentHash, TypeHelloSyn),
		BodyHelloSyn{
			Persona: persona,
		},
	})
}

func (msg *HelloSyn) GetBody() Body {
	return msg.Body
}

func (msg *HelloSyn) Check(vault *types.Vault, state *types.State) error {
	auth := state.GetAuth()
	if !auth.IsPublic() && !auth.CanRead(msg.GetClientAddr()) {
		return code.NonReadPermission
	}
	return nil
}

func (msg *HelloSyn) Execute(vault *types.Vault, state *types.State) error {
	err := state.Join(msg.Body.Persona)
	if err != nil {
		return err
	}
	// secretKey := vault.GetSecretKey()
	// encryptedSecretKey, err := msg.Body.Persona.PubKey.Encrypt(secretKey.Bytes())
	// if err != nil {
	// 	return err
	// }

	// peerID := state.GetPeerID()
	// clientAddr := vault.GetAddr()
	// personae := state.GetPersonae()
	// auth := state.GetAuth()
	// msgAck := NewMsgHelloAck(peerID, clientAddr, Hash(msg), personae, auth, encryptedSecretKey)
	// err = box.Publish(msgAck, false)
	// if err != nil {
	// 	return err
	// }
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

func NewMsgHelloAck(peerID types.ID, clientAddr crypto.Addr, parentHash types.Hash,
	personae types.Personae, auth *types.Auth, encryptedSecretKey []byte) *Msg {
	return NewMsg(&HelloAck{
		NewHead(peerID, clientAddr, parentHash, TypeHelloAck),
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

func (msg HelloAck) Check(vault *types.Vault, state *types.State) error {
	// parentMsg, err := msg.getParentMsg(box)
	// if err != nil {
	// 	return err
	// }
	// if parentMsg == nil {
	// 	return code.NonExistingParent
	// }
	return nil
}

func (msg HelloAck) Execute(vault *types.Vault, state *types.State) error {
	privKey := vault.GetPrivKey()
	secretKeyByte, err := privKey.Decrypt(msg.Body.EncryptedSecretKey)
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
	if util.HasField("personae", state) {
		state.SetPersonae(msg.Body.Personae)
	}
	if util.HasField("auth", state) {
		state.SetAuth(msg.Body.Auth)
	}
	if util.HasField("secretKey", vault) {
		vault.SetSecretKey(secretKey)
	}
	return nil
}
