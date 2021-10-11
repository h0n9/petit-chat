package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivKeyMarshal(t *testing.T) {
	privKey, err := GenPrivKey()
	assert.Nil(t, err)

	b, err := privKey.MarshalJSON()
	assert.Nil(t, err)

	newPrivKey := PrivKey{}
	err = newPrivKey.UnmarshalJSON(b)
	assert.Nil(t, err)

	assert.EqualValues(t, privKey, newPrivKey)
}
func TestPubKeyMarshal(t *testing.T) {
	privKey, err := GenPrivKey()
	assert.Nil(t, err)

	pubKey := privKey.PubKey()

	b, err := pubKey.MarshalJSON()
	assert.Nil(t, err)

	newPubKey := PubKey{}
	err = newPubKey.UnmarshalJSON(b)
	assert.Nil(t, err)

	assert.EqualValues(t, pubKey, newPubKey)
}
