package msg

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type BodyUpdate struct {
	Auth     *types.Auth    `json:"auth"`
	Personae types.Personae `json:"personae"`
}

type Update struct {
	Head
	Body BodyUpdate `json:"body"`
}

func NewMsgUpdate(peerID types.ID, clientAddr crypto.Addr, parentHash types.Hash, auth *types.Auth, personae types.Personae) *Msg {
	return NewMsg(&Update{
		NewHead(peerID, clientAddr, parentHash, TypeUpdate),
		BodyUpdate{
			Auth:     auth,
			Personae: personae,
		},
	})
}

func (msg *Update) GetBody() Body {
	return msg.Body
}

func (msg *Update) Check(vault *types.Vault, state *types.State) error {
	auth := state.GetAuth()
	if !auth.CanExecute(msg.GetClientAddr()) {
		return code.NonExecutePermission
	}
	return nil
}
func (msg *Update) Execute(vault *types.Vault, state *types.State) error {
	if util.HasField("auth", state) {
		state.SetAuth(msg.Body.Auth)
	}
	if util.HasField("personae", state) {
		state.SetPersonae(msg.Body.Personae)
	}
	return nil
}
