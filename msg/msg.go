package msg

import (
	"encoding/json"
	"time"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type Body interface {
	Check(*Box, crypto.Addr) error
	Execute(*Box, types.Hash) error
}

type Signature struct {
	PubKey   crypto.PubKey `json:"pubkey"`
	SigBytes []byte        `json:"sig_bytes"`
}

type Msg struct {
	Hash       types.Hash `json:"hash"`
	Timestamp  time.Time  `json:"timestamp"`
	PeerID     types.ID   `json:"peer_id"`
	ParentHash types.Hash `json:"parent_hash"`
	Body       Body       `json:"body"`
	Signature  Signature  `json:"signature"`
}

type MsgToSign struct {
	Hash       types.Hash `json:"-"`
	Timestamp  time.Time  `json:"timestamp"`
	PeerID     types.ID   `json:"peer_id"`
	ParentHash types.Hash `json:"parent_hash"`
	Body       Body       `json:"body"`
	Signature  Signature  `json:"-"`
}

type MsgToVerify MsgToSign

type MsgCapsule struct {
	Encrypted bool   `json:"encrpyted"`
	Type      Type   `json:"type"`
	Data      []byte `json:"data"`
}

type MsgEx struct {
	Read     bool `json:"read"`
	Received bool `json:"received"`
	*Msg
}

func NewMsg(peerID types.ID, parentHash types.Hash, body Body) *Msg {
	return &Msg{
		Timestamp:  time.Now(),
		PeerID:     peerID,
		ParentHash: parentHash,
		Body:       body,
	}
}

func (msg *Msg) GetPeerID() types.ID {
	return msg.PeerID
}

func (msg *Msg) GetBody() Body {
	return msg.Body
}

func (msg *Msg) SetBody(body Body) {
	msg.Body = body
}

func (msg *Msg) GetTimestamp() time.Time {
	return msg.Timestamp
}

func (msg *Msg) GetHash() types.Hash {
	return msg.Hash
}

func (msg *Msg) GetSignature() Signature {
	return msg.Signature
}

func (msg *Msg) GetParentHash() types.Hash {
	return msg.ParentHash
}

func (msg *Msg) IsEOS() bool {
	switch msg.Body.(type) {
	case *BodyBye:
		return true
	}
	return false
}

func (msg *Msg) Sign(privKey *crypto.PrivKey) error {
	pubKey := privKey.PubKey()
	msgToSign := MsgToSign(*msg)
	b, err := json.Marshal(msgToSign)
	if err != nil {
		return err
	}
	sigBytes, err := privKey.Sign(b)
	if err != nil {
		return err
	}

	msg.Hash = util.ToSHA256(b)
	msg.Signature = Signature{
		SigBytes: sigBytes,
		PubKey:   pubKey,
	}

	return nil
}

func (msg *Msg) Verify() error {
	msgToVerify := MsgToVerify(*msg)
	b, err := json.Marshal(msgToVerify)
	if err != nil {
		return err
	}
	sigBytes := msg.GetSignature().SigBytes
	ok := msg.Signature.PubKey.Verify(b, sigBytes)
	if !ok {
		return code.FailedToVerify
	}

	return nil
}

func MarshalJSON(msg *Msg) ([]byte, error) {
	return json.Marshal(msg)
}

func UnmarshalJSON(data []byte, msg *Msg) error {
	return json.Unmarshal(data, msg)
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
