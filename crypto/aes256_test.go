package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAES256EncryptDecrypt(t *testing.T) {
	sk1, err := GenSecretKey()
	assert.Nil(t, err)
	sk2, err := GenSecretKey()
	assert.Nil(t, err)

	msg := []byte("hello world let's make the world better")

	cipher, err := sk1.Encrypt(msg)
	assert.Nil(t, err)

	assert.NotEqual(t, msg, cipher)

	plain, err := sk1.Decrypt(cipher)
	assert.Nil(t, err)

	assert.Equal(t, msg, plain)

	_, err = sk2.Decrypt(cipher)
	assert.NotNil(t, err)
}
