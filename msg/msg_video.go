package msg

import (
	"encoding/json"
)

type MsgStructVideo struct {
	Data      []byte `json:"data"`
	Extension string `json:"extension"`
}

func NewMsgStructVideo(data []byte, extension string) *MsgStructVideo {
	return &MsgStructVideo{Data: data, Extension: extension}
}

func (msv *MsgStructVideo) Encapsulate() ([]byte, error) {
	return json.Marshal(msv)
}

func (msv *MsgStructVideo) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msv)
}
