package msg

import (
	"encoding/json"

	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type MsgStructUpdateSyn struct {
	Auth     *types.Auth    `json:"auth"`
	Personae types.Personae `json:"personae"`
}

func NewMsgStructUpdateSyn(auth *types.Auth, personae types.Personae) *MsgStructUpdateSyn {
	return &MsgStructUpdateSyn{Auth: auth, Personae: personae}
}

func (msu *MsgStructUpdateSyn) Encapsulate() ([]byte, error) {
	return json.Marshal(msu)
}

func (msu *MsgStructUpdateSyn) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msu)
}

func (msu *MsgStructUpdateSyn) Execute(b *Box) error {
	if util.HasField("auth", b) {
		b.auth = msu.Auth
	}
	if util.HasField("personae", b) {
		b.personae = msu.Personae
	}

	return nil
}
