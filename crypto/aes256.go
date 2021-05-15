package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

type SecretKey struct {
	key    []byte
	block  cipher.Block
	aesGCM cipher.AEAD
	nonce  []byte
}

func NewSecretKey(key []byte) (*SecretKey, error) {
	// check constraint
	if len(key) != SecretKeySize {
		return nil, fmt.Errorf("wrong size for SecretKey: %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())

	return &SecretKey{
		key:    key,
		block:  block,
		aesGCM: aesGCM,
		nonce:  nonce,
	}, nil
}

func GenSecretKey() (*SecretKey, error) {
	key := make([]byte, SecretKeySize)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}

	return NewSecretKey(key)
}

func (sk *SecretKey) Encrypt(plaintext []byte) ([]byte, error) {
	return sk.aesGCM.Seal(nil, sk.nonce, plaintext, nil), nil
}

func (sk *SecretKey) Decrypt(ciphertext []byte) ([]byte, error) {
	return sk.aesGCM.Open(nil, sk.nonce, ciphertext, nil)
}
