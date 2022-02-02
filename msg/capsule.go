package msg

import (
	"encoding/json"
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

func (mc *MsgCapsule) Bytes() ([]byte, error) {
	return json.Marshal(mc)
}
