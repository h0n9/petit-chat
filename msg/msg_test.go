package msg

import (
	"fmt"
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

type MsgTest struct {
	Head
	Body BodyTest `json:"body"`
}

func (msg *MsgTest) GetBody() Body {
	return msg.Body
}

func (msg *MsgTest) Check(box *Box) error {
	return nil
}

func (msg *MsgTest) Execute(box *Box) error {
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

	box := Box{
		vault: types.NewVault(nil, privKey1, nil),
	}

	msg := NewMsg(&MsgTest{
		Head{
			Timestamp:  time.Now(),
			PeerID:     types.ID(""),
			ParentHash: types.Hash{},
		},
		BodyTest{
			Name:    "nothing",
			Content: "this is nothing.",
		},
	})

	fmt.Println("before:", msg)

	err = box.Sign(msg)
	assert.Nil(t, err)

	fmt.Println("after:", msg)

	err = box.Verify(msg)
	assert.Nil(t, err)

	// manipulate pubKey on purpose
	msg.Signature.PubKey = pubKey2

	err = box.Verify(msg)
	assert.Equal(t, err, code.FailedToVerify)
}
