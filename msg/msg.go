package msg

import (
	"encoding/json"
	"time"

	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type Msg struct {
	Timestamp time.Time `json:"timestamp"`
	From      types.ID  `json:"from"` // always ONE from
	Data      []byte    `json:"data"`
}

type MsgEx struct {
	Read     bool `json:"read"`
	Received bool `json:"received"`
	*Msg
}

func NewMsg(from types.ID, data []byte) *Msg {
	return &Msg{
		Timestamp: time.Now(),
		From:      from,
		Data:      data,
	}
}

func (msg *Msg) GetFrom() types.ID {
	return msg.From
}

func (msg *Msg) GetData() []byte {
	return msg.Data
}

func (msg *Msg) GetTime() time.Time {
	return msg.Timestamp
}

func (msg *Msg) Hash() (types.Hash, error) {
	b, err := msg.MarshalJSON()
	if err != nil {
		return types.Hash{}, err
	}
	return util.ToSHA256(b), nil
}

func (msg *Msg) MarshalJSON() ([]byte, error) {
	return json.Marshal(*msg)
}

func UnmarshalJSON(data []byte) (*Msg, error) {
	msg := Msg{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}
