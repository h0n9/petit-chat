package msg

import (
	"time"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
)

type Body interface{}

type Signature struct {
	PubKey   crypto.PubKey `json:"pubkey"`
	SigBytes []byte        `json:"sig_bytes"`
}

type Msg struct {
	Hash      types.Hash `json:"hash"`
	Signature Signature  `json:"signature"`
	Base      `json:"base"`
	Metas     types.Metas `json:"-"`
}

func NewMsg(base Base) *Msg {
	return &Msg{
		Base:  base,
		Metas: make(types.Metas),
	}
}

func (msg *Msg) GetHash() types.Hash {
	return msg.Hash
}

func (msg *Msg) SetHash(hash types.Hash) {
	msg.Hash = hash
}

func (msg *Msg) GetSignature() Signature {
	return msg.Signature
}

func (msg *Msg) SetSignature(signature Signature) {
	msg.Signature = signature
}

func (msg *Msg) GetMetas() types.Metas {
	return msg.Metas
}

func (msg *Msg) SetMetas(metas types.Metas) {
	msg.Metas = metas
}

func (msg *Msg) GetMeta(addr crypto.Addr) types.Meta {
	return msg.Metas[addr]
}

func (msg *Msg) SetMeta(addr crypto.Addr, meta types.Meta) {
	msg.Metas[addr] = meta
}

func (msg *Msg) MergeMeta(addr crypto.Addr, newMeta types.Meta) {
	oldMeta, exist := msg.Metas[addr]
	if exist {
		newMeta |= oldMeta
	}
	msg.SetMeta(addr, newMeta)
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

	// ops
	check(*Box) error
	Check(*Box) error
	Execute(*Box) error
}

type Head struct {
	Timestamp  time.Time   `json:"timestamp"`
	PeerID     types.ID    `json:"peer_id"`
	ClientAddr crypto.Addr `json:"client_addr"`
	ParentHash types.Hash  `json:"parent_hash"`
	Type       Type        `json:"type"`
}

func NewHead(box *Box, parentHash types.Hash, msgType Type) Head {
	return Head{
		Timestamp:  time.Now(),
		PeerID:     box.GetMyID(),
		ClientAddr: box.GetMyPersona().Address,
		ParentHash: parentHash,
		Type:       msgType,
	}
}

type MsgCapsule struct {
	Encrypted bool   `json:"encrpyted"`
	Type      Type   `json:"type"`
	Data      []byte `json:"data"`
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

func (msg *Head) getParentMsg(b *Box) (*Msg, error) {
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

func (msg *Head) check(b *Box) error {
	// check msg.ParentMsgHash
	pm, err := msg.getParentMsg(b)
	if err != nil {
		return err
	}
	if pm != nil && !pm.GetParentHash().IsEmpty() {
		return code.AlreadyHavingParent
	}

	// TODO: add more constraints
	return nil
}
