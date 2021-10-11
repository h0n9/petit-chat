package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: TestSignVerify
func TestECDSASignVerify(t *testing.T) {
	privKey, err := GenPrivKey()
	assert.Nil(t, err)

	pubKey := privKey.PubKey()

	msg1 := []byte("hello world let's make the better world")
	msg2 := []byte("hello world let's make the world better")

	sigBytes, err := privKey.Sign(msg1)
	assert.Nil(t, err)

	ok := pubKey.Verify(msg1, sigBytes)
	assert.True(t, ok)

	ok = pubKey.Verify(msg2, sigBytes)
	assert.False(t, ok)
}
