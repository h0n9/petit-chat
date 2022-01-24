package types

import (
	"github.com/h0n9/petit-chat/crypto"
)

type Vault struct {
	id        ID
	persona   *Persona // TODO: move to client side permanently
	privKey   *crypto.PrivKey
	secretKey *crypto.SecretKey
}

func NewVault(id ID, persona *Persona, privKey *crypto.PrivKey, secretKey *crypto.SecretKey) *Vault {
	return &Vault{
		id:        id,
		persona:   persona,
		privKey:   privKey,
		secretKey: secretKey,
	}
}

func (v *Vault) GetID() ID {
	return v.id
}

func (v *Vault) GetPersona() *Persona {
	return v.persona
}

func (v *Vault) GetPrivKey() *crypto.PrivKey {
	return v.privKey
}

func (v *Vault) GetPubKey() crypto.PubKey {
	return v.privKey.PubKey()
}

func (v *Vault) GetAddr() crypto.Addr {
	return v.persona.Address
}

func (v *Vault) GetSecretKey() *crypto.SecretKey {
	return v.secretKey
}

func (v *Vault) SetSecretKey(secretKey *crypto.SecretKey) {
	v.secretKey = secretKey
}
