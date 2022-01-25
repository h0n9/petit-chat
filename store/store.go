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

func (s *Store) getNextIndex() (types.Index, []byte, error) {
	indexByteSlice, err := s.rootDB.Get(nextIndexKey)
	if err != nil {
		return 0, []byte{}, err
	}
	index, err := types.IndexFromByteSlice(indexByteSlice)
	if err != nil {
		return 0, []byte{}, err
	}
	return index, indexByteSlice, nil
}

func (s *Store) GetNextIndex() (types.Index, error) {
	index, _, err := s.getNextIndex()
	return index, err

}

func (s *Store) Append(hash types.Hash, data []byte) (types.Index, error) {
	hashByteSlice := hash[:]
	index, indexByteSlice, err := s.getNextIndex()
	if err != nil {
		return types.Index(0), err
	}
	err = s.indexDB.Set(indexByteSlice, hashByteSlice)
	if err != nil {
		return types.Index(0), err
	}
	err = s.hashDB.Set(hashByteSlice, data)
	if err != nil {
		return types.Index(0), err
	}
	err = s.SetNextIndex(index + types.Index(1))
	if err != nil {
		return types.Index(0), err
	}
	return index, nil
}

func (s *Store) getDataByHash(hash []byte) ([]byte, error) {
	return s.hashDB.Get(hash)
}

func (s *Store) GetDataByHash(hash types.Hash) ([]byte, error) {
	return s.getDataByHash(hash[:])
}

func (s *Store) GetDataByIndex(index types.Index) ([]byte, error) {
	indexByteSlice := types.IndexToByteSlice(index)
	hashByteSlice, err := s.indexDB.Get(indexByteSlice)
	if err != nil {
		return nil, err
	}
	return s.getDataByHash(hashByteSlice)
}
