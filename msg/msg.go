package msg

import (
	"encoding/json"
	"time"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type Msg struct {
	Timestamp     time.Time  `json:"timestamp"`
	From          types.ID   `json:"from"` // always ONE from
	Type          MsgType    `json:"type"`
	ParentMsgHash types.Hash `json:"parent_msg_hash"`
	Data          []byte     `json:"data"`
}

type MsgEx struct {
	Read     bool `json:"read"`
	Received bool `json:"received"`
	*Msg
}

func NewMsg(from types.ID, msgType MsgType, parentMsgHash types.Hash, data []byte) *Msg {
	return &Msg{
		Timestamp:     time.Now(),
		From:          from,
		Type:          msgType,
		ParentMsgHash: parentMsgHash,
		Data:          data,
	}
}

func (msg *Msg) GetFrom() types.ID {
	return msg.From
}

func (msg *Msg) GetType() MsgType {
	return msg.Type
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

func (msg *Msg) IsEOS() bool {
	return msg.Type == MsgTypeEOS
}

func (msg *Msg) Encapsulate() ([]byte, error) {
	// TODO: change to other format (later)
	return msg.MarshalJSON()
}

func Decapsulate(data []byte) (*Msg, error) {
	// TODO: change to other format (later)
	return UnmarshalJSON(data)
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

func (msg *Msg) getParentMsg(b *Box) (*Msg, error) {
	// check if parentMsgHash is empty
	pmh := msg.ParentMsgHash
	if types.IsEmpty(pmh) {
		return nil, nil
	}
	// get msg corresponding to msgHash
	pm := b.GetMsg(pmh)
	if pm == nil {
		return nil, code.NonExistingParentMsg
	}
	return pm, nil
}
