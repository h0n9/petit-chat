package chat

import (
	"github.com/h0n9/petit-chat/msg"
	"github.com/h0n9/petit-chat/p2p"
	"github.com/h0n9/petit-chat/util"
)

var (
	node     *p2p.Node
	hostPeer msg.Peer
)

func NewChatCmd(n *p2p.Node, h msg.Peer) *util.Cmd {
	node = n
	hostPeer = h
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
