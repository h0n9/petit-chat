package types

import (
	"testing"

	"github.com/h0n9/petit-chat/crypto"
	"github.com/stretchr/testify/assert"
)

func genRandomAddr() (Addr, error) {
	privKey, err := crypto.GenPrivKey()
	if err != nil {
		return Addr(""), err
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

	perms := map[Addr]Perm{
		addrTest1: permNone,
		addrTest2: permRead,
		addrTest3: permWrite,
		addrTest4: permExecute,
	}

	a := NewAuth(false, perms)

	// IsPublic()
	public := a.IsPublic()
	assert.False(t, public)

	// CheckMinPerm()
	_, err = a.CheckMinPerm(addrTest0)
	assert.Error(t, err)

	ok, err := a.CheckMinPerm(addrTest1)
	assert.NoError(t, err)
	assert.False(t, ok)

	ok, err = a.CheckMinPerm(addrTest2)
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = a.CheckMinPerm(addrTest3)
	assert.NoError(t, err)
	assert.False(t, ok)

	ok, err = a.CheckMinPerm(addrTest4)
	assert.NoError(t, err)
	assert.False(t, ok)

	// add addrTest0
	err = a.SetPerm(addrTest0, permNone)
	assert.NoError(t, err)
	permTest0, err := a.getPerm(addrTest0)
	assert.NoError(t, err)
	assert.NotEqual(t, permTest0, permRead)
	assert.Equal(t, permTest0, permNone)

	ok, err = a.CheckMinPerm(addrTest0)
	assert.NoError(t, err)
	assert.False(t, ok)

	err = a.SetPerm(addrTest0, permRead)
	assert.NoError(t, err)
	permTest0, err = a.getPerm(addrTest0)
	assert.NoError(t, err)
	assert.NotEqual(t, permTest0, permNone)
	assert.Equal(t, permTest0, permRead)

	ok, err = a.CheckMinPerm(addrTest0)
	assert.NoError(t, err)
	assert.True(t, ok)

	err = a.SetPerm(addrTest0, permRead|permWrite)
	assert.NoError(t, err)
	permTest0, err = a.getPerm(addrTest0)
	assert.NoError(t, err)
	assert.NotEqual(t, permTest0, permRead)
	assert.Equal(t, permTest0, permRead|permWrite)

	ok, err = a.CheckMinPerm(addrTest0)
	assert.NoError(t, err)
	assert.True(t, ok)
}
