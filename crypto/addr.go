package crypto

import (
	"bytes"
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

func (addr Addr) IsDrivenFrom(pubKey PubKey) bool {
	addrB := pubKey.Address()
	return addr.Equals(addrB)
}

func (addr Addr) MarshalJSON() ([]byte, error) {
	data := make([]byte, len(addr)+2)
	data[0] = '"'
	data[len(data)-1] = '"'
	copy(data[1:], addr.String())
	return data, nil
}

func (addr *Addr) UnmarshalJSON(data []byte) error {
	copy(addr[:], data[1:len(data)-1])
	return nil
}

func (addr Addr) Equals(addrB Addr) bool {
	return bytes.Equal(addr.Bytes(), addrB.Bytes())
}

func (addr Addr) Bytes() []byte {
	return addr[:]
}

func (addr Addr) String() string {
	return string(addr[:])
}
