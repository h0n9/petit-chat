package msg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
)

type BodyTest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (body *BodyTest) Check(box *Box, hash types.Hash, addr crypto.Addr) error {
	return nil
}

func (body *BodyTest) Execute(box *Box, hash types.Hash, addr crypto.Addr) error {
	return nil
}

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
		PeerID:     types.ID(""),
		ParentHash: types.Hash{},
		Body: &BodyTest{
			Name:    "nothing",
			Content: "this is nothing.",
		},
	}
	err = msg.Sign(&privKey1)
	assert.Nil(t, err)

	err = msg.Verify()
	assert.Nil(t, err)

	// manipulate pubKey on purpose
	msg.Signature.PubKey = pubKey2

	err = msg.Verify()
	assert.Equal(t, err, code.FailedToVerify)
}
