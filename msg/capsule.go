package msg

import (
	"encoding/json"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
)

type MsgCapsule struct {
	Encrypted bool   `json:"encrpyted"`
	Type      Type   `json:"type"`
	Data      []byte `json:"data"`
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
	return nil
}

func (mc *MsgCapsule) Decapsulate(secretKey *crypto.SecretKey) (*Msg, error) {
	data := mc.Data
	if mc.Encrypted {
		decryptedData, err := secretKey.Decrypt(mc.Data)
		if err != nil {
			return nil, err
		}
		data = decryptedData
	}

	m := NewMsg(mc.Type.Base())
	if m == nil {
		return nil, code.UnknownMsgType
	}

	err := json.Unmarshal(data, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (mc *MsgCapsule) Bytes() ([]byte, error) {
	return json.Marshal(mc)
}
