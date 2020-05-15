package peer

import (
	"github.com/h0n9/petit-chat/p2p"
	"github.com/h0n9/petit-chat/util"
)

var node *p2p.Node

func NewPeerCmd(n *p2p.Node) *util.Cmd {
	node = n
	return util.NewCmd(
		"peer",
		"peer related commands",
		nil,
		listCmd,
		addCmd,
		removeCmd,
		blackCmd,
	)
}
