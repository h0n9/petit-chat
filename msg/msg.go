package msg

import (
	"encoding/json"
	"time"

	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type Body interface{}

type Msg struct {
	Base  `json:"base"`
	Metas types.Metas `json:"-"`
}

func NewMsg(base Base) *Msg {
	return &Msg{
		Base:  base,
		Metas: make(types.Metas),
	}
}

func (msg *Msg) encode() ([]byte, error) {
	return json.Marshal(msg)
}

func (msg *Msg) Encapsulate() (*Capsule, error) {
	data, err := msg.encode()
	if err != nil {
		return nil, err
	}
	return NewCapsule(false, msg.GetType(), data), nil
}

func (msg *Msg) Hash() (types.Hash, error) {
	data, err := msg.encode()
	if err != nil {
		return types.EmptyHash, err
	}
	return util.ToSHA256(data), nil
}

type Base interface {
	// accessors
	GetTimestamp() time.Time
	GetPeerID() types.ID
	GetClientAddr() crypto.Addr
	GetParentHash() types.Hash
	GetType() Type
	GetBody() Body
	IsEOS() bool

	check() error

	// ops
	Check(types.Hash, Helper) error
	Execute(types.Hash, Helper) error
}

type Head struct {
	Timestamp  time.Time   `json:"timestamp"`
	PeerID     types.ID    `json:"peer_id"`
	ClientAddr crypto.Addr `json:"client_addr"`
	ParentHash types.Hash  `json:"parent_hash"`
	Type       Type        `json:"type"`
}

func NewHead(peerID types.ID, clientAddr crypto.Addr, parentHash types.Hash, msgType Type) Head {
	return Head{
		Timestamp:  time.Now(),
		PeerID:     peerID,
		ClientAddr: clientAddr,
		ParentHash: parentHash,
		Type:       msgType,
	}
}

func (msg *Head) GetTimestamp() time.Time {
	return msg.Timestamp
}

func (msg *Head) GetPeerID() types.ID {
	return msg.PeerID
}

func (msg *Head) GetClientAddr() crypto.Addr {
	return msg.ClientAddr
}

func (msg *Head) GetParentHash() types.Hash {
	return msg.ParentHash
}

func (msg *Head) GetType() Type {
	return msg.Type
}

func (msg *Head) IsEOS() bool {
	return msg.Type == TypeBye
}

func (msg *Head) check() error {
	err := msg.Type.Check()
	if err != nil {
		return err
	}
	// TODO: add more constraints
	return nil
}

type Helper interface {
	// accessors
	GetVault() *types.Vault
	GetState() *types.State
	GetStore() *CapsuleStore
	GetPeerID() types.ID

	// operators
	Publish(msg *Msg, encrypt bool) error
}
