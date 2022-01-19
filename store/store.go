package store

import (
	tmdb "github.com/tendermint/tm-db"

	"github.com/h0n9/petit-chat/types"
)

var (
	nextIndexKey = []byte("next_index")
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
	err := store.SetNextIndex(types.Index(0))
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (s *Store) SetNextIndex(index types.Index) error {
	return s.rootDB.Set(nextIndexKey, types.IndexToByteSlice(index))
}

func (s *Store) GetNextIndex() (types.Index, error) {
	value, err := s.rootDB.Get(nextIndexKey)
	if err != nil {
		return 0, err
	}
	return types.IndexFromByteSlice(value)
}
