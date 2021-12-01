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

func (msus *MsgStructUpdateSyn) Encapsulate() ([]byte, error) {
	return json.Marshal(msus)
}

func (msus *MsgStructUpdateSyn) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msus)
}

func (msus *MsgStructUpdateSyn) Execute(b *Box) error {
	if util.HasField("auth", b) {
		b.auth = msus.Auth
	}
	if util.HasField("personae", b) {
		b.personae = msus.Personae
	}
	return nil
}
