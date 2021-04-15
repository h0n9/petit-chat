package msg

import (
	"encoding/json"
	"time"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type MsgFrom struct {
	PeerID     types.ID    `json:"peer_id"`
	ClientAddr crypto.Addr `json:"client_addr"`
}

type Msg struct {
	Timestamp     time.Time  `json:"timestamp"`
	From          MsgFrom    `json:"from"`
	Type          types.Msg  `json:"type"`
	ParentMsgHash types.Hash `json:"parent_msg_hash"`
	Data          []byte     `json:"data"`
}

type MsgEx struct {
	Read     bool `json:"read"`
	Received bool `json:"received"`
	*Msg
}

func NewMsg(pID types.ID, cAddr crypto.Addr, msgType types.Msg, parentMsgHash types.Hash, data []byte) *Msg {
	return &Msg{
		Timestamp: time.Now(),
		From: MsgFrom{
			PeerID:     pID,
			ClientAddr: cAddr,
		},
		Type:          msgType,
		ParentMsgHash: parentMsgHash,
		Data:          data,
	}
}

func (msg *Msg) GetFrom() MsgFrom {
	return msg.From
}

func (msg *Msg) GetType() types.Msg {
	return msg.Type
}

func (msg *Msg) GetData() []byte {
	return msg.Data
}

func (msg *Msg) GetTime() time.Time {
	return msg.Timestamp
}

func (msg *Msg) Hash() (types.Hash, error) {
	b, err := MarshalJSON(msg)
	if err != nil {
		return types.Hash{}, err
	}
	return util.ToSHA256(b), nil
}

func (msg *Msg) IsEOS() bool {
	return msg.Type == types.MsgBye
}

func (msg *Msg) Encapsulate() ([]byte, error) {
	// TODO: change to other format (later)
	return MarshalJSON(msg)
}

func (msg *Msg) Decapsulate(data []byte) error {
	// TODO: change to other format (later)
	return UnmarshalJSON(data, msg)
}

func MarshalJSON(msg *Msg) ([]byte, error) {
	return json.Marshal(*msg)
}

func UnmarshalJSON(data []byte, msg *Msg) error {
	return json.Unmarshal(data, msg)
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
