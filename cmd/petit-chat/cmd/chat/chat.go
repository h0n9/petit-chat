package chat

import (
	"github.com/h0n9/petit-chat/p2p"
	"github.com/h0n9/petit-chat/util"
)

var node *p2p.Node

func NewChatCmd(n *p2p.Node) *util.Cmd {
	node = n
	return util.NewCmd(
		"chat",
		"chat related commands",
		nil,
		listCmd,
		enterCmd,
		createCmd,
		leaveCmd,
	)
}
