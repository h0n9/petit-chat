package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestECIESEncryptDecrypt(t *testing.T) {
	priv1, err := GenPrivKey()
	assert.Nil(t, err)

	priv2, err := GenPrivKey()
	assert.Nil(t, err)

	msg := []byte("hello world let's make the better world")

	cipher, err := priv2.PubKey().Encrypt(msg)
	assert.Nil(t, err)

	plain, err := priv2.Decrypt(cipher)
	assert.Nil(t, err)

	assert.Equal(t, msg, plain)

	_, err = priv1.Decrypt(cipher)
	assert.NotNil(t, err)
}
