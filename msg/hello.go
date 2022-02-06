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

func (msg *HelloSyn) Check(hash types.Hash, helper Helper) error {
	state := helper.GetState()

	auth := state.GetAuth()
	if !auth.IsPublic() && !auth.CanRead(msg.GetClientAddr()) {
		return code.NonReadPermission
	}
	return nil
}

func (msg *HelloSyn) Execute(hash types.Hash, helper Helper) error {
	vault := helper.GetVault()
	state := helper.GetState()
	peerID := helper.GetPeerID()

	err := state.Join(msg.Body.Persona)
	if err != nil {
		return err
	}
	secretKey := vault.GetSecretKey()
	encryptedSecretKey, err := msg.Body.Persona.PubKey.Encrypt(secretKey.Bytes())
	if err != nil {
		return err
	}

	clientAddr := vault.GetAddr()
	personae := state.GetPersonae()
	auth := state.GetAuth()
	msgAck := NewMsgHelloAck(peerID, clientAddr, hash, personae, auth, encryptedSecretKey)
	err = helper.Publish(msgAck, false)
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

func (msg *HelloAck) Check(hash types.Hash, helper Helper) error {
	store := helper.GetStore()

	pmh := msg.GetParentHash()
	pc := store.GetCapsule(pmh)
	if pc == nil {
		return code.NonExistingParent
	}
	return nil
}

func (msg *HelloAck) Execute(hash types.Hash, helper Helper) error {
	vault := helper.GetVault()
	state := helper.GetState()

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
