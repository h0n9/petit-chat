package msg

import (
	"encoding/json"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type (
	PubSub = pubsub.PubSub
	Sub    = pubsub.Subscription
)

type Msg struct {
	read     bool
	received bool

	Timestamp time.Time `json:"timestamp"`
	From      Peer      `json:"from"` // always ONE from
	To        []Peer    `json:"to"`   // could be SEVERAL to
	Data      []byte    `json:"data"`
}

func NewMsg(data []byte, from Peer, to []Peer) *Msg {
	return &Msg{
		read:     false,
		received: false,

		Timestamp: time.Now(),
		From:      from,
		To:        to,
		Data:      data,
	}
}

func (msg *Msg) GetFrom() Peer {
	return msg.From
}

func (msg *Msg) GetData() []byte {
	return msg.Data
}

func (msg *Msg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg)
}

func UnmarshalJSON(data []byte) (*Msg, error) {
	msg := Msg{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}
