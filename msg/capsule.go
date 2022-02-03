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

type MsgCapsule struct {
	Hash      types.Hash `json:"hash"`
	Signature Signature  `json:"signature"`
	Type      Type       `json:"type"`
	Data      []byte     `json:"data"`

	Encrypted bool `json:"encrpyted"`
}

func NewMsgCapsule(encrypted bool, msgType Type, data []byte) *MsgCapsule {
	return &MsgCapsule{
		Encrypted: encrypted,
		Type:      msgType,
		Data:      data,
	}
}

func NewMsgCapsuleFromBytes(data []byte) (*MsgCapsule, error) {
	msgCapsule := MsgCapsule{}
	err := json.Unmarshal(data, &msgCapsule)
	if err != nil {
		return nil, err
	}
	return &msgCapsule, nil
}

func (mc *MsgCapsule) Check() error {
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

func (mc *MsgCapsule) GetHash() types.Hash {
	return mc.Hash
}

func (mc *MsgCapsule) SetHash(hash types.Hash) {
	mc.Hash = hash
}

func (mc *MsgCapsule) GetSignature() Signature {
	return mc.Signature
}

func (mc *MsgCapsule) SetSignature(signature Signature) {
	mc.Signature = signature
}

func (mc *MsgCapsule) Sign(privKey *crypto.PrivKey) error {
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

func (mc *MsgCapsule) Verify() error {
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

func (mc *MsgCapsule) Encrypt(secretKey *crypto.SecretKey) error {
	encryptedData, err := secretKey.Encrypt(mc.Data)
	if err != nil {
		return err
	}
	mc.Data = encryptedData
	mc.Encrypted = true
	return nil
}

func (mc *MsgCapsule) Decrypt(secretKey *crypto.SecretKey) error {
	decryptedData, err := secretKey.Decrypt(mc.Data)
	if err != nil {
		return err
	}
	mc.Data = decryptedData
	mc.Encrypted = false
	return nil
}

func (mc *MsgCapsule) Decapsulate() (*Msg, error) {
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

func (mc *MsgCapsule) Bytes() ([]byte, error) {
	return json.Marshal(mc)
}
