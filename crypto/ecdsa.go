package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	h "crypto/sha256"
	"math/big"
)

func (privKey PrivKey) Sign(msg []byte) ([]byte, error) {
	hash := h.Sum256(msg)

	priv := privKey.ToECDSA()
	r, s, err := ecdsa.Sign(rand.Reader, priv, hash[:])
	if err != nil {
		return nil, err
	}

	sig := make([]byte, SigSize)
	rb := r.Bytes()
	sb := s.Bytes()

	copy(sig[32-len(rb):], rb)
	copy(sig[64-len(sb):], sb)

	return sig, nil
}

func (pubKey PubKey) Verify(msg []byte, sig []byte) bool {
	if len(sig) != SigSize {
		return false
	}
	hash := h.Sum256(msg)
	return ecdsa.Verify(
		pubKey.ToECDSA(),
		hash[:],
		new(big.Int).SetBytes(sig[:32]),
		new(big.Int).SetBytes(sig[32:]),
	)
}
