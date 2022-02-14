package types

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/h0n9/petit-chat/crypto"
)

func TestStateMeta(t *testing.T) {
	state := NewState(true)
	assert.NotNil(t, state)

	hash := Hash{0x0}
	addr := crypto.Addr("addr")
	meta := NewMeta(true, false, false)

	metas, exist := state.GetMetas(hash)
	assert.Equal(t, false, exist)
	assert.Equal(t, Metas(nil), metas)

	state.UpdateMeta(hash, addr, meta)

	metas, exist = state.GetMetas(hash)
	assert.Equal(t, true, exist)
	assert.Equal(t, Metas{addr: meta}, metas)

	meta_0 := NewMeta(true, true, false)
	state.UpdateMeta(hash, addr, meta_0)

	metas, exist = state.GetMetas(hash)
	assert.Equal(t, true, exist)
	assert.Equal(t, Metas{addr: meta_0}, metas)

	// meta := NewMeta(false, true, false)
	// state.UpdateMeta(hash, addr, meta)
	//
	// metas, exist = state.GetMetas(hash)
	// assert.Equal(t, true, exist)
	// assert.Equal(t, Metas{addr: meta}, metas)

	hash = Hash{0x1}
	meta_1 := NewMeta(true, false, false)
	state.UpdateMeta(hash, addr, meta_1)

	metas, exist = state.GetMetas(Hash{0x0})
	assert.Equal(t, true, exist)
	assert.Equal(t, Metas{addr: meta_0}, metas)
	metas, exist = state.GetMetas(Hash{0x1})
	assert.Equal(t, true, exist)
	assert.Equal(t, Metas{addr: meta_1}, metas)
}
