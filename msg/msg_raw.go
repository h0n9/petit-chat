package msg

import (
	"encoding/json"
)

type MsgStructRaw struct {
	Data []byte `json:"data"`
}

func NewMsgStructRaw(data []byte) *MsgStructRaw {
	return &MsgStructRaw{Data: data}
}

func (msr *MsgStructRaw) GetData() []byte {
	return msr.Data
}

func (msr *MsgStructRaw) Encapsulate() ([]byte, error) {
	return json.Marshal(msr)
}

func (msr *MsgStructRaw) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msr)
}
