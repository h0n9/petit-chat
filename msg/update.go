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

func (body *BodyUpdate) Check(box *Box, from *From) error {
	if !box.auth.CanExecute(from.ClientAddr) {
		return code.NonExecutePermission
	}
	return nil
}
func (body *BodyUpdate) Execute(box *Box, hash types.Hash) error {
	if util.HasField("auth", box) {
		box.auth = body.Auth
	}
	if util.HasField("personae", box) {
		box.personae = body.Personae
	}
	return nil
}
