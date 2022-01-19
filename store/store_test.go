package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tmdb "github.com/tendermint/tm-db"

	"github.com/h0n9/petit-chat/types"
)

func TestNextIndex(t *testing.T) {
	s, err := NewStore(tmdb.NewMemDB())
	assert.NoError(t, err)

	// non-equal
	err = s.SetNextIndex(0)
	assert.NoError(t, err)
	latestIndex, err := s.GetNextIndex()
	assert.NoError(t, err)
	assert.NotEqual(t, types.Index(1), latestIndex)

	// equals
	// zero
	err = s.SetNextIndex(0)
	assert.NoError(t, err)
	latestIndex, err = s.GetNextIndex()
	assert.NoError(t, err)
	assert.Equal(t, types.Index(0), latestIndex)

	// biggest number in uint64
	err = s.SetNextIndex(18446744073709551615)
	assert.NoError(t, err)
	latestIndex, err = s.GetNextIndex()
	assert.NoError(t, err)
	assert.Equal(t, types.Index(18446744073709551615), latestIndex)
}
