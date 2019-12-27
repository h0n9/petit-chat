package crypto

import (
	"crypto/elliptic"
	h "crypto/sha256"

	"github.com/btcsuite/btcutil/base58"
)

var preAddr = []byte{'p', 'c', 'h'}

func (pubKey PubKey) Address() Addr {
	addr := Addr{}

	pubKeyByte := elliptic.Marshal(c, pubKey.X(), pubKey.Y())
	hash := h.Sum256(pubKeyByte)

	b58 := base58.Encode(hash[:])

	// addr := preAddr + base58.Encode(hash(pubKey))
	copy(addr[0:], preAddr)
	copy(addr[len(preAddr):], b58)

	return addr
}

func (addr Addr) String() string {
	return string(addr[:])
}
