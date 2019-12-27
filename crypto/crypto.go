package crypto

import (
	"crypto/elliptic"
)

/*
	- cryptographic algorithm to use
	  : ecdsa + p256
	- key type
	  : PrivKey, PubKey
	  (base) ecdsa.PrivateKey, ecdsa.PublicKey
	- address type
	  : PRE_ADDR + base58.Encode(hash(pubKey))
*/

var (
	c = elliptic.P256()
)

const (
	PrivKeySize = 32
	PubKeySize  = 65
	AddrSize    = 47
	SigSize     = 64
)

type (
	PrivKey [PrivKeySize]byte
	PubKey  [PubKeySize]byte
	Addr    [AddrSize]byte
	Hash    [32]byte
)
