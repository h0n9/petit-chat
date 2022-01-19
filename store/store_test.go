package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tmdb "github.com/tendermint/tm-db"

	"github.com/h0n9/petit-chat/types"
)

func TestLatestCount(t *testing.T) {
	s := NewStore(tmdb.NewMemDB())

	// non-equal
	err := s.SetLastestCount(types.Count(0))
	assert.NoError(t, err)
	latestCount, err := s.GetLatestCount()
	assert.NoError(t, err)
	assert.NotEqual(t, types.Count(1), latestCount)

	// equals
	// zero
	err = s.SetLastestCount(types.Count(0))
	assert.NoError(t, err)
	latestCount, err = s.GetLatestCount()
	assert.NoError(t, err)
	assert.Equal(t, types.Count(0), latestCount)

	// biggest number in uint64
	err = s.SetLastestCount(types.Count(18446744073709551615))
	assert.NoError(t, err)
	latestCount, err = s.GetLatestCount()
	assert.NoError(t, err)
	assert.Equal(t, types.Count(18446744073709551615), latestCount)
}
