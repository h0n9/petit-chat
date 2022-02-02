package types

import (
	"github.com/h0n9/petit-chat/crypto"
)

type Vault struct {
	persona   *Persona
	privKey   *crypto.PrivKey
	secretKey *crypto.SecretKey
}

func NewVault(persona *Persona, privKey *crypto.PrivKey, secretKey *crypto.SecretKey) *Vault {
	return &Vault{
		persona:   persona,
		privKey:   privKey,
		secretKey: secretKey,
	}
}

func (v *Vault) GetPersona() *Persona {
	return v.persona
}

func (v *Vault) GetPrivKey() *crypto.PrivKey {
	return v.privKey
}

func (v *Vault) GetPubKey() *crypto.PubKey {
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
