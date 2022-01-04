package types

import (
	"bytes"
	"encoding/hex"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

const (
	hashSize = 32
)

var (
	EmptyHash = Hash{}
)

type (
	ID        = peer.ID
	Host      = host.Host
	PubSub    = pubsub.PubSub
	Sub       = pubsub.Subscription
	Topic     = pubsub.Topic
	PubSubMsg = pubsub.Message
	Hash      [hashSize]byte
)

func (hash Hash) String() string {
	return string(hash[:])
}

func (hash *Hash) MarshalJSON() ([]byte, error) {
	return []byte(`"` + hex.EncodeToString(hash[:]) + `"`), nil
}

func (hash *Hash) UnmarshalJSON(data []byte) error {
	tmp, err := hex.DecodeString(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}
	copy(hash[:], tmp)
	return nil
}

func (hash Hash) IsEmpty() bool {
	return bytes.Equal(hash[:], EmptyHash[:])
}
