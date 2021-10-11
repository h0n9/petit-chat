package msg

import (
	"encoding/json"
)

type MsgStructAudio struct {
	Data      []byte `json:"data"`
	Extension string `json:"extension"`
}

func NewMsgStructAudio(data []byte, extension string) *MsgStructAudio {
	return &MsgStructAudio{Data: data, Extension: extension}
}

func (msa *MsgStructAudio) Encapsulate() ([]byte, error) {
	return json.Marshal(msa)
}

func (msa *MsgStructAudio) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msa)
}
