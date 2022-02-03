package msg

import (
	"encoding/json"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type Signature struct {
	PubKey   *crypto.PubKey `json:"pubkey"`
	SigBytes []byte         `json:"sig_bytes"`
}

type Capsule struct {
	Hash      types.Hash `json:"hash"`
	Signature Signature  `json:"signature"`
	Type      Type       `json:"type"`
	Data      []byte     `json:"data"`

	Encrypted bool `json:"encrpyted"`
}

func NewCapsule(encrypted bool, msgType Type, data []byte) *Capsule {
	return &Capsule{
		Encrypted: encrypted,
		Type:      msgType,
		Data:      data,
	}
}

func NewCapsuleFromBytes(data []byte) (*Capsule, error) {
	capsule := Capsule{}
	err := json.Unmarshal(data, &capsule)
	if err != nil {
		return nil, err
	}
	return &capsule, nil
}

func (mc *Capsule) Check() error {
	err := mc.Type.Check()
	if err != nil {
		return err
	}
	if !mc.Encrypted {
		err = mc.Verify()
		if err != nil {
			return err
		}
	}
	return nil
}

func (mc *Capsule) GetHash() types.Hash {
	return mc.Hash
}

func (mc *Capsule) SetHash(hash types.Hash) {
	mc.Hash = hash
}

func (mc *Capsule) GetSignature() Signature {
	return mc.Signature
}

func (mc *Capsule) SetSignature(signature Signature) {
	mc.Signature = signature
}

func (mc *Capsule) Sign(privKey *crypto.PrivKey) error {
	hash := util.ToSHA256(mc.Data)
	sigBytes, err := privKey.Sign(mc.Data)
	if err != nil {
		return err
	}
	mc.SetHash(hash)
	mc.SetSignature(Signature{
		SigBytes: sigBytes,
		PubKey:   privKey.PubKey(),
	})
	return nil
}

func (mc *Capsule) Verify() error {
	signature := mc.GetSignature()
	if signature.SigBytes == nil {
		return code.ImproperSigBytes
	}
	if signature.PubKey == nil {
		return code.ImproperPubKey
	}
	ok := signature.PubKey.Verify(mc.Data, signature.SigBytes)
	if !ok {
		return code.FailedToVerify
	}
	return nil
}

func (mc *Capsule) Encrypt(secretKey *crypto.SecretKey) error {
	encryptedData, err := secretKey.Encrypt(mc.Data)
	if err != nil {
		return err
	}
	mc.Data = encryptedData
	mc.Encrypted = true
	return nil
}

func (mc *Capsule) Decrypt(secretKey *crypto.SecretKey) error {
	decryptedData, err := secretKey.Decrypt(mc.Data)
	if err != nil {
		return err
	}
	mc.Data = decryptedData
	mc.Encrypted = false
	return nil
}

func (mc *Capsule) Decapsulate() (*Msg, error) {
	m := NewMsg(mc.Type.Base())
	if m == nil {
		return nil, code.UnknownMsgType
	}

	err := json.Unmarshal(mc.Data, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (mc *Capsule) Bytes() ([]byte, error) {
	return json.Marshal(mc)
}

type CapsuleStore struct {
	length        types.Length
	capsules      []*Capsule              // TODO: limit the size of msgs slice
	capsuleHashes map[types.Hash]*Capsule // TODO: limit the size of msgHashes map
}

func NewCapsuleStore() *CapsuleStore {
	return &CapsuleStore{
		length:        types.Length(0),
		capsules:      make([]*Capsule, 0),
		capsuleHashes: make(map[types.Hash]*Capsule),
	}
}

func (cs *CapsuleStore) Append(capsule *Capsule) (types.Index, error) {
	hash := capsule.GetHash()
	if cs.Has(hash) {
		return 0, code.AlreadyAppendedCapsule
	}

	cs.capsules = append(cs.capsules, capsule)
	cs.capsuleHashes[hash] = capsule

	index := cs.length
	cs.length += 1

	return index, nil
}

func (cs *CapsuleStore) GetCapsules() []*Capsule {
	return cs.capsules
}

func (cs *CapsuleStore) GetCapsule(hash types.Hash) *Capsule {
	return cs.capsuleHashes[hash]
}

func (cs *CapsuleStore) Has(hash types.Hash) bool {
	_, exist := cs.capsuleHashes[hash]
	return exist
}
