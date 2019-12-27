package crypto

import (
	"crypto/rand"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

func (pubKey PubKey) Encrypt(msg []byte) ([]byte, error) {
	pubKeyECDSA := pubKey.ToECDSA()
	pubKeyECIES := ecies.ImportECDSAPublic(pubKeyECDSA)

	ct, err := ecies.Encrypt(rand.Reader, pubKeyECIES, msg, nil, nil)
	if err != nil {
		return nil, err
	}

	return ct, nil
}

func (privKey PrivKey) Decrypt(cipher []byte) ([]byte, error) {
	privKeyECDSA := privKey.ToECDSA()
	privKeyECIES := ecies.ImportECDSA(privKeyECDSA)

	msg, err := privKeyECIES.Decrypt(cipher, nil, nil)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
