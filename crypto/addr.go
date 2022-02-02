package crypto

import (
	"crypto/elliptic"
	h "crypto/sha256"

	"github.com/btcsuite/btcutil/base58"
)

var preAddr = []byte{'p', 'c', 'h'}

func (pubKey PubKey) Address() Addr {
	tmp := make([]byte, AddrSize)
	pubKeyByte := elliptic.Marshal(c, pubKey.X(), pubKey.Y())
	hash := h.Sum256(pubKeyByte)
	b58 := base58.Encode(hash[:])
	// addr := preAddr + base58.Encode(hash(pubKey))
	copy(tmp[0:], preAddr)
	copy(tmp[len(preAddr):], b58)
	return Addr(tmp)
}

func (addrA Addr) IsDrivenFrom(pubKey *PubKey) bool {
	addrB := pubKey.Address()
	return addrA.Equals(addrB)
}

func (addrA Addr) Equals(addrB Addr) bool {
	return addrA == addrB
}
