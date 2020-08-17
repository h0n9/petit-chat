package msg

import (
	"encoding/json"
	"time"

	"github.com/h0n9/petit-chat/types"
)

type Msg struct {
	read     bool
	received bool

	timestamp time.Time `json:"timestamp"`
	from      types.ID  `json:"from"` // always ONE from
	data      []byte    `json:"data"`
}

func NewMsg(from types.ID, data []byte) *Msg {
	return &Msg{
		read:     false,
		received: false,

		timestamp: time.Now(),
		from:      from,
		data:      data,
	}
}

func (msg *Msg) GetFrom() types.ID {
	return msg.from
}

func (msg *Msg) GetData() []byte {
	return msg.data
}

func (msg *Msg) GetTime() time.Time {
	return msg.timestamp
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
