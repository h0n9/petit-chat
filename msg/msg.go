package msg

import (
	"encoding/json"
	"time"

	"github.com/h0n9/petit-chat/types"
)

type Msg struct {
	read     bool
	received bool

	Timestamp time.Time  `json:"timestamp"`
	From      types.ID   `json:"from"` // always ONE from
	To        []types.ID `json:"to"`   // could be SEVERAL to
	Data      []byte     `json:"data"`
}

func NewMsg(data []byte, from types.ID, to []types.ID) *Msg {
	return &Msg{
		read:     false,
		received: false,

		Timestamp: time.Now(),
		From:      from,
		To:        to,
		Data:      data,
	}
}

func (msg *Msg) GetFrom() types.ID {
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
