package msg

import (
	"encoding/json"

	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type MsgStructUpdate struct {
	Auth     *types.Auth    `json:"auth"`
	Personae types.Personae `json:"personae"`
}

func NewMsgStructUpdate(auth *types.Auth, personae types.Personae) *MsgStructUpdate {
	return &MsgStructUpdate{Auth: auth, Personae: personae}
}

func (msu *MsgStructUpdate) Encapsulate() ([]byte, error) {
	return json.Marshal(msu)
}

func (msu *MsgStructUpdate) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msu)
}

func (msu *MsgStructUpdate) Execute(b *Box) error {
	if util.HasField("auth", b) {
		b.auth = msu.Auth
	}
	if util.HasField("personae", b) {
		b.personae = msu.Personae
	}

	return nil
}
