package msg

import (
	"encoding/json"

	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type MsgStructUpdateAck struct {
	Auth     *types.Auth    `json:"auth"`
	Personae types.Personae `json:"personae"`
}

func NewMsgStructUpdateAck(auth *types.Auth, personae types.Personae) *MsgStructUpdateAck {
	return &MsgStructUpdateAck{Auth: auth, Personae: personae}
}

func (msua *MsgStructUpdateAck) Encapsulate() ([]byte, error) {
	return json.Marshal(msua)
}

func (msua *MsgStructUpdateAck) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msua)
}

func (msua *MsgStructUpdateAck) Execute(b *Box) error {
	if util.HasField("auth", b) {
		b.auth = msua.Auth
	}
	if util.HasField("personae", b) {
		b.personae = msua.Personae
	}
	return nil
}
