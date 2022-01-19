package store

import (
	tmdb "github.com/tendermint/tm-db"

	"github.com/h0n9/petit-chat/msg"
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

func (s *Store) Append(msg msg.Msg) error {
	hash := msg.GetHash()
	hashByteSlice := hash[:]
	msgByteSlice, err := msg.ToByteSlice()
	if err != nil {
		return err
	}
	nextIndex, nextIndexByteSlice, err := s.getNextIndex()
	if err != nil {
		return err
	}
	err = s.indexDB.Set(nextIndexByteSlice, hashByteSlice)
	if err != nil {
		return err
	}
	err = s.hashDB.Set(hashByteSlice, msgByteSlice)
	if err != nil {
		return err
	}
	err = s.SetNextIndex(nextIndex + types.Index(1))
	if err != nil {
		return err
	}
	return nil
}
