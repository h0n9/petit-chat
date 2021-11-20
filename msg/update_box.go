package msg

import (
	"encoding/json"

	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type MsgStructUpdateBox struct {
	Auth     *types.Auth    `json:"auth"`
	Personae types.Personae `json:"personae"`
}

func NewMsgStructUpdateBox(auth *types.Auth, personae types.Personae) *MsgStructUpdateBox {
	return &MsgStructUpdateBox{Auth: auth, Personae: personae}
}

func (msub *MsgStructUpdateBox) Encapsulate() ([]byte, error) {
	return json.Marshal(msub)
}

func (msub *MsgStructUpdateBox) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msub)
}

func (msub *MsgStructUpdateBox) Execute(b *Box) error {
	if util.HasField("auth", b) {
		b.auth = msub.Auth
	}
	if util.HasField("personae", b) {
		b.personae = msub.Personae
	}

	return nil
}
