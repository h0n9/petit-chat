package types

import (
	"testing"

	"github.com/h0n9/petit-chat/crypto"
	"github.com/stretchr/testify/assert"
)

func genRandomAddr() (crypto.Addr, error) {
	privKey, err := crypto.GenPrivKey()
	if err != nil {
		return crypto.Addr(""), err
	}

	return privKey.PubKey().Address(), nil

}

func TestAuth(t *testing.T) {
	addrTest0, err := genRandomAddr()
	assert.NoError(t, err)
	addrTest1, err := genRandomAddr()
	assert.NoError(t, err)
	addrTest2, err := genRandomAddr()
	assert.NoError(t, err)
	addrTest3, err := genRandomAddr()
	assert.NoError(t, err)
	addrTest4, err := genRandomAddr()
	assert.NoError(t, err)

	perms := map[crypto.Addr]Perm{
		addrTest1: permNone,
		addrTest2: permRead,
		addrTest3: permWrite,
		addrTest4: permExecute,
	}

	a := NewAuth(false, perms)

	// IsPublic()
	public := a.IsPublic()
	assert.False(t, public)

	// add addrTest0
	err = a.SetPerm(addrTest0, permNone)
	assert.NoError(t, err)
	permTest0, err := a.getPerm(addrTest0)
	assert.NoError(t, err)
	assert.NotEqual(t, permTest0, permRead)
	assert.Equal(t, permTest0, permNone)

	err = a.SetPerm(addrTest0, permRead)
	assert.NoError(t, err)
	permTest0, err = a.getPerm(addrTest0)
	assert.NoError(t, err)
	assert.NotEqual(t, permTest0, permNone)
	assert.Equal(t, permTest0, permRead)

	err = a.SetPerm(addrTest0, permRead|permWrite)
	assert.NoError(t, err)
	permTest0, err = a.getPerm(addrTest0)
	assert.NoError(t, err)
	assert.NotEqual(t, permTest0, permRead)
	assert.Equal(t, permTest0, permRead|permWrite)
}
