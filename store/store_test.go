package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tmdb "github.com/tendermint/tm-db"

	"github.com/h0n9/petit-chat/msg"
	"github.com/h0n9/petit-chat/types"
)

func TestNextIndex(t *testing.T) {
	s, err := NewStore(tmdb.NewMemDB())
	assert.NoError(t, err)

	// check if successfully initialized with 0
	nextIndex, err := s.GetNextIndex()
	assert.NoError(t, err)
	assert.NotEqual(t, types.Index(1), nextIndex)

	// non-equal
	err = s.SetNextIndex(1)
	assert.NoError(t, err)
	nextIndex, err = s.GetNextIndex()
	assert.NoError(t, err)
	assert.NotEqual(t, types.Index(0), nextIndex)

	// equals
	// zero
	err = s.SetNextIndex(0)
	assert.NoError(t, err)
	nextIndex, err = s.GetNextIndex()
	assert.NoError(t, err)
	assert.Equal(t, types.Index(0), nextIndex)

	// biggest number in uint64
	err = s.SetNextIndex(18446744073709551615)
	assert.NoError(t, err)
	nextIndex, err = s.GetNextIndex()
	assert.NoError(t, err)
	assert.Equal(t, types.Index(18446744073709551615), nextIndex)
}

type BodyTest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type MsgTest struct {
	msg.Head
	Body BodyTest `json:"body"`
}

func (msg *MsgTest) GetBody() msg.Body {
	return msg.Body
}

func (msg *MsgTest) Check(hash types.Hash, helper msg.Helper) error {
	return nil
}

func (msg *MsgTest) Execute(hash types.Hash, helper msg.Helper) error {
	return nil
}

// func TestAppend(t *testing.T) {
// 	// prepare msg to append
// 	m := msg.NewMsg(&MsgTest{
// 		msg.Head{
// 			Timestamp:  time.Now(),
// 			PeerID:     types.ID(""),
// 			ParentHash: types.Hash{},
// 		},
// 		BodyTest{
// 			Name:    "nothing",
// 			Content: "this is nothing.",
// 		},
// 	})
// 	hash := msg.Hash(m)
// 	m.SetHash(hash)
// 	mData, err := json.Marshal(m)
// 	assert.NoError(t, err)
//
// 	s, err := NewStore(tmdb.NewMemDB())
// 	assert.NoError(t, err)
//
// 	nextIndex, err := s.GetNextIndex()
// 	assert.NoError(t, err)
// 	assert.Equal(t, types.Index(0), nextIndex)
//
// 	index, err := s.Append(hash, mData)
// 	assert.NoError(t, err)
// 	assert.Equal(t, types.Index(0), index)
//
// 	nextIndex, err = s.GetNextIndex()
// 	assert.NoError(t, err)
// 	assert.Equal(t, types.Index(1), nextIndex)
//
// 	mmData, err := s.GetDataByHash(hash)
// 	assert.NoError(t, err)
// 	assert.EqualValues(t, mData, mmData)
//
// 	mmData, err = s.GetDataByIndex(index)
// 	assert.NoError(t, err)
// 	assert.EqualValues(t, mData, mmData)
// }
