package peer

import (
	"github.com/h0n9/petit-chat/msg"
	"github.com/h0n9/petit-chat/p2p"
	"github.com/h0n9/petit-chat/util"
)

var (
	node     *p2p.Node
	hostPeer msg.Peer
)

func NewPeerCmd(n *p2p.Node, h msg.Peer) *util.Cmd {
	node = n
	hostPeer = h
	return util.NewCmd(
		"peer",
		"peer related commands",
		nil,
		listCmd,
		addCmd,
		removeCmd,
		blockCmd,
	)
}
