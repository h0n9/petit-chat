package msg

import (
	"github.com/h0n9/petit-chat/code"
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

func (msg *Update) GetBody() Body {
	return msg.Body
}

func (msg *Update) Check(box *Box) error {
	if !box.auth.CanExecute(msg.GetClientAddr()) {
		return code.NonExecutePermission
	}
	return nil
}
func (msg *Update) Execute(box *Box) error {
	if util.HasField("auth", box) {
		box.auth = msg.Body.Auth
	}
	if util.HasField("personae", box) {
		box.personae = msg.Body.Personae
	}
	return nil
}
