package types

import (
	"testing"

	"github.com/h0n9/petit-chat/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
)

func genRandomID() (ID, error) {
	id := ID("")

	pk, err := crypto.GenPrivKey()
	if err != nil {
		return id, err
	}

	pkECDSAP2P, err := pk.ToECDSAP2P()
	if err != nil {
		return id, err
	}

	newId, err := peer.IDFromPrivateKey(pkECDSAP2P)
	if err != nil {
		return id, err
	}

	id = newId

	return id, nil

}

func TestAuth(t *testing.T) {
	idTest0, err := genRandomID()
	assert.NoError(t, err)
	idTest1, err := genRandomID()
	assert.NoError(t, err)
	idTest2, err := genRandomID()
	assert.NoError(t, err)
	idTest3, err := genRandomID()
	assert.NoError(t, err)
	idTest4, err := genRandomID()
	assert.NoError(t, err)

	perms := map[ID]Perm{
		idTest1: permNone,
		idTest2: permRead,
		idTest3: permWrite,
		idTest4: permExecute,
	}

	a := NewAuth(false, perms)

	// IsPublic()
	public := a.IsPublic()
	assert.False(t, public)

	// CheckMinPerm()
	_, err = a.CheckMinPerm(idTest0)
	assert.Error(t, err)

	ok, err := a.CheckMinPerm(idTest1)
	assert.NoError(t, err)
	assert.False(t, ok)

	ok, err = a.CheckMinPerm(idTest2)
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = a.CheckMinPerm(idTest3)
	assert.NoError(t, err)
	assert.False(t, ok)

	ok, err = a.CheckMinPerm(idTest4)
	assert.NoError(t, err)
	assert.False(t, ok)

	// add idTest0
	err = a.SetPerm(idTest0, permNone)
	assert.NoError(t, err)
	permTest0, err := a.GetPerm(idTest0)
	assert.NoError(t, err)
	assert.NotEqual(t, permTest0, permRead)
	assert.Equal(t, permTest0, permNone)

	ok, err = a.CheckMinPerm(idTest0)
	assert.NoError(t, err)
	assert.False(t, ok)

	err = a.SetPerm(idTest0, permRead)
	assert.NoError(t, err)
	permTest0, err = a.GetPerm(idTest0)
	assert.NoError(t, err)
	assert.NotEqual(t, permTest0, permNone)
	assert.Equal(t, permTest0, permRead)

	ok, err = a.CheckMinPerm(idTest0)
	assert.NoError(t, err)
	assert.True(t, ok)

	err = a.SetPerm(idTest0, permRead|permWrite)
	assert.NoError(t, err)
	permTest0, err = a.GetPerm(idTest0)
	assert.NoError(t, err)
	assert.NotEqual(t, permTest0, permRead)
	assert.Equal(t, permTest0, permRead|permWrite)

	ok, err = a.CheckMinPerm(idTest0)
	assert.NoError(t, err)
	assert.True(t, ok)
}
