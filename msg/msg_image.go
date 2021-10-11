package msg

import (
	"encoding/json"
)

type MsgStructImage struct {
	Data      []byte `json:"data"`
	Extension string `json:"extension"`
}

func NewMsgStructImage(data []byte, extension string) *MsgStructImage {
	return &MsgStructImage{Data: data, Extension: extension}
}

func (msi *MsgStructImage) Encapsulate() ([]byte, error) {
	return json.Marshal(msi)
}

func (msi *MsgStructImage) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msi)
}
