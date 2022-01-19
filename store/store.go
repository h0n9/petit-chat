package store

import (
	tmdb "github.com/tendermint/tm-db"

	"github.com/h0n9/petit-chat/types"
)

var (
	latestCountKey = []byte("latest_count")
	prefixIndex    = []byte("index:")
	prefixHash     = []byte("hash:")
)

type Store struct {
	rootDB  tmdb.DB
	indexDB tmdb.DB
	hashDB  tmdb.DB
}

func NewStore(rootDB tmdb.DB) *Store {
	store := Store{
		rootDB:  rootDB,
		indexDB: tmdb.NewPrefixDB(rootDB, prefixIndex),
		hashDB:  tmdb.NewPrefixDB(rootDB, prefixHash),
	}
	return &store
}

func (s *Store) SetLastestCount(count types.Count) error {
	return s.rootDB.Set(latestCountKey, types.CountToByteSlice(count))
}

func (s *Store) GetLatestCount() (types.Count, error) {
	value, err := s.rootDB.Get(latestCountKey)
	if err != nil {
		return 0, err
	}
	return types.CountFromByteSlice(value)
}
