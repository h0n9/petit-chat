package msg

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/types"
)

type BodyMeta struct {
	Meta types.Meta `json:"meta"`
}

type Meta struct {
	Head
	Body BodyMeta `json:"body"`
}

func (msg *Meta) GetBody() Body {
	return msg.Body
}

func (msg *Meta) Check(box *Box) error {
	clientAddr := msg.GetClientAddr()
	if msg.Body.Meta.Received() || msg.Body.Meta.Read() {
		if !box.auth.IsPublic() && !box.auth.CanRead(clientAddr) {
			return code.NonReadPermission
		}
	}
	if msg.Body.Meta.Typing() && !box.auth.CanWrite(clientAddr) {
		return code.NonWritePermission
	}
	return nil
}

func (msg *Meta) Execute(box *Box) error {
	parentMsg, err := msg.getParentMsg(box)
	if err != nil {
		return err
	}
	parentMsg.MergeMeta(msg.GetClientAddr(), msg.Body.Meta)
	return nil
}
