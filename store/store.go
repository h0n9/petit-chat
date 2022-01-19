package store

import (
	tmdb "github.com/tendermint/tm-db"

	"github.com/h0n9/petit-chat/types"
)

var (
	nextCountKey = []byte("next_count")
	prefixIndex  = []byte("index:")
	prefixHash   = []byte("hash:")
)

type Store struct {
	rootDB  tmdb.DB
	indexDB tmdb.DB
	hashDB  tmdb.DB
}

func NewStore(rootDB tmdb.DB) (*Store, error) {
	store := Store{
		rootDB:  rootDB,
		indexDB: tmdb.NewPrefixDB(rootDB, prefixIndex),
		hashDB:  tmdb.NewPrefixDB(rootDB, prefixHash),
	}
	err := store.SetNextCount(types.Count(0))
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (s *Store) SetNextCount(count types.Count) error {
	return s.rootDB.Set(nextCountKey, types.CountToByteSlice(count))
}

func (s *Store) GetNextCount() (types.Count, error) {
	value, err := s.rootDB.Get(nextCountKey)
	if err != nil {
		return 0, err
	}
	return types.CountFromByteSlice(value)
}
