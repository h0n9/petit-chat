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
	Timestamp  time.Time  `json:"timestamp"`
	From       MsgFrom    `json:"from"`
	Type       MsgType    `json:"type"`
	ParentHash types.Hash `json:"parent_hash"`
	Encrypted  bool       `json:"encrypted"`
	Data       []byte     `json:"data"`
	Hash       types.Hash
}

type MsgEx struct {
	Read     bool `json:"read"`
	Received bool `json:"received"`
	*Msg
}

func NewMsg(pID types.ID, cAddr crypto.Addr,
	msgType MsgType, parentHash types.Hash, encrypted bool, data []byte,
) (*Msg, error) {
	msg := Msg{
		Timestamp: time.Now(),
		From: MsgFrom{
			PeerID:     pID,
			ClientAddr: cAddr,
		},
		Type:       msgType,
		ParentHash: parentHash,
		Encrypted:  encrypted,
		Data:       data,
	}
	hash, err := msg.hash()
	if err != nil {
		return nil, err
	}
	msg.Hash = hash
	return &msg, nil
}

func (msg *Msg) GetFrom() MsgFrom {
	return msg.From
}

func (msg *Msg) GetType() MsgType {
	return msg.Type
}

func (msg *Msg) GetData() []byte {
	return msg.Data
}

func (msg *Msg) SetData(data []byte) {
	msg.Data = data
}

func (msg *Msg) GetTimestamp() time.Time {
	return msg.Timestamp
}

func (msg *Msg) GetHash() types.Hash {
	return msg.Hash
}

func (msg *Msg) GetParentHash() types.Hash {
	return msg.ParentHash
}

func (msg *Msg) IsEOS() bool {
	return msg.Type == MsgTypeBye
}

func (msg *Msg) Encapsulate() ([]byte, error) {
	// TODO: change to other format (later)
	return MarshalJSON(msg)
}

func (msg *Msg) Decapsulate(data []byte) error {
	// TODO: change to other format (later)
	err := UnmarshalJSON(data, msg)
	if err != nil {
		return err
	}
	hash, err := msg.hash()
	if err != nil {
		return err
	}
	msg.Hash = hash
	return nil
}

func MarshalJSON(msg *Msg) ([]byte, error) {
	return json.Marshal(*msg)
}

func UnmarshalJSON(data []byte, msg *Msg) error {
	return json.Unmarshal(data, msg)
}

func (msg *Msg) hash() (types.Hash, error) {
	b, err := MarshalJSON(msg)
	if err != nil {
		return types.Hash{}, err
	}
	return util.ToSHA256(b), nil
}

func (msg *Msg) getParentMsg(b *Box) (*Msg, error) {
	// check if parentMsgHash is empty
	pmh := msg.ParentHash
	if pmh.IsEmpty() {
		return nil, nil
	}
	// get msg corresponding to msgHash
	pm := b.GetMsg(pmh)
	if pm == nil {
		// TODO: this error should be optional
		return nil, code.NonExistingParent
	}
	return pm, nil
}

func (msg *Msg) check(b *Box) error {
	// check msgType
	mt := msg.GetType()
	err := mt.Check()
	if err != nil {
		return err
	}

	// check msg.ParentMsgHash
	pm, err := msg.getParentMsg(b)
	if err != nil {
		return err
	}
	if pm != nil && !pm.ParentHash.IsEmpty() {
		return code.AlreadyHavingParent
	}

	// TODO: add more constraints

	return nil
}
