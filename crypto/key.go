package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/h0n9/petit-chat/code"
	lc "github.com/libp2p/go-libp2p-core/crypto"
)

func GenPrivKey() (PrivKey, error) {
	privKey := PrivKey{}

	p256, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		return privKey, err
	}

	copy(privKey[:], p256.D.Bytes())
	return privKey, nil
}

func (privKey PrivKey) Bytes() []byte {
	return privKey[:]
}

func (privKey PrivKey) Equals(target PrivKey) bool {
	return bytes.Equal(privKey.Bytes(), target.Bytes())
}

func (privKey PrivKey) String() string {
	return hex.EncodeToString(privKey[:])
}

func (privKey PrivKey) MarshalJSON() ([]byte, error) {
	data := make([]byte, PrivKeySize*2+2)
	data[0] = '"'
	data[len(data)-1] = '"'
	copy(data[1:], privKey.String())
	return data, nil
}

func (privKey *PrivKey) UnmarshalJSON(data []byte) error {
	if len(data) != PrivKeySize*2+2 {
		return fmt.Errorf("privKeyJSON size %d != expected %d",
			len(data), PrivKeySize*2+2,
		)
	}

	_, err := hex.Decode(privKey[:], data[1:len(data)-1])
	if err != nil {
		return err
	}

	return nil
}

func (privKey PrivKey) ToECDSA() *ecdsa.PrivateKey {
	X, Y := c.ScalarBaseMult(privKey[:])
	return &ecdsa.PrivateKey{
		D: new(big.Int).SetBytes(privKey[:]),
		PublicKey: ecdsa.PublicKey{
			Curve: c,
			X:     X,
			Y:     Y,
		},
	}
}

func (privKey PrivKey) ToECDSAP2P() (lc.PrivKey, error) {
	pk, _, err := lc.ECDSAKeyPairFromKey(privKey.ToECDSA())
	if err != nil {
		return nil, err
	}

	return pk, nil
}

/*
func (privKey PrivKey) FromECDSA(*ecdsa.PrivateKey) {
	copy(privKey[:], *ecdsa.PrivateKey.D.Bytes())
}
*/

// PubKey related functions

func (privKey PrivKey) PubKey() PubKey {
	pubKey := PubKey{PubKeyPrefix}

	priv := privKey.ToECDSA()
	X := priv.X.Bytes()
	Y := priv.Y.Bytes()

	copy(pubKey[33-len(X):], X)
	copy(pubKey[65-len(Y):], Y)

	return pubKey
}

func (pubKey PubKey) Check() error {
	if len(pubKey) != PubKeySize {
		return code.ImproperPubKey
	}
	if pubKey[0] != PubKeyPrefix {
		return code.ImproperPubKey
	}
	return nil
}

func (pubKey PubKey) Bytes() []byte {
	return pubKey[:]
}

func (pubKey PubKey) ToECDSA() *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: c,
		X:     new(big.Int).SetBytes(pubKey[1:33]),
		Y:     new(big.Int).SetBytes(pubKey[33:]),
	}
}

func (pubKey PubKey) Equals(target PubKey) bool {
	return bytes.Equal(pubKey.Bytes(), target.Bytes())
}

func (pubKey PubKey) String() string {
	return hex.EncodeToString(pubKey[:])
}

func (pubKey PubKey) MarshalJSON() ([]byte, error) {
	data := make([]byte, PubKeySize*2+2)
	data[0] = '"'
	data[len(data)-1] = '"'

	copy(data[1:], pubKey.String())
	return data, nil
}

func (pubKey PubKey) UnmarshalJSON(data []byte) error {
	if len(data) != PubKeySize*2+2 {
		return fmt.Errorf("pubKeyJSON size %d != expected %d",
			len(data), PubKeySize*2+2,
		)
	}

	_, err := hex.Decode(pubKey[:], data[1:len(data)-1])
	if err != nil {
		return err
	}

	return nil
}

func (pubKey PubKey) X() *big.Int {
	return pubKey.ToECDSA().X
}

func (pubKey PubKey) Y() *big.Int {
	return pubKey.ToECDSA().Y
}
