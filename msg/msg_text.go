package msg

import (
	"encoding/json"
)

type MsgStructText struct {
	Data     []byte `json:"data"`
	Encoding string `json:"encoding"`
}

func NewMsgStructText(data []byte, encoding string) *MsgStructText {
	return &MsgStructText{Data: data, Encoding: encoding}
}

func (mst *MsgStructText) GetData() []byte {
	return mst.Data
}

func (mst *MsgStructText) Encapsulate() ([]byte, error) {
	return json.Marshal(mst)
}

func (mst *MsgStructText) Decapsulate(data []byte) error {
	return json.Unmarshal(data, mst)
}
