package msg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
)

func TestMsgSignVerify(t *testing.T) {
	privKey1, err := crypto.GenPrivKey()
	assert.Nil(t, err)
	privKey2, err := crypto.GenPrivKey()
	assert.Nil(t, err)

	pubKey1 := privKey1.PubKey()
	assert.NotNil(t, pubKey1)
	pubKey2 := privKey2.PubKey()
	assert.NotNil(t, pubKey2)

	msg := Msg{
		Timestamp:  time.Now(),
		From:       From{},
		Type:       MsgTypeRaw,
		ParentHash: types.Hash{},
		Encrypted:  false,
		Data:       []byte("hello world"),
	}
	err = msg.Sign(&privKey1)
	assert.Nil(t, err)

	err = msg.Verify(&pubKey1)
	assert.Nil(t, err)
	err = msg.Verify(&pubKey2)
	assert.Equal(t, err, code.FailedToVerify)
}
