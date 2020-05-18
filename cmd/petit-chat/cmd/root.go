package cmd

import (
	"github.com/h0n9/petit-chat/cmd/petit-chat/cmd/chat"
	"github.com/h0n9/petit-chat/cmd/petit-chat/cmd/peer"
	"github.com/h0n9/petit-chat/p2p"
	"github.com/h0n9/petit-chat/util"
)

var node *p2p.Node

func NewRootCmd(n *p2p.Node) *util.Cmd {
	node = n
	return util.NewCmd("petit-chat", "entry point for petit-chat", nil,
		infoCmd,
		peer.NewPeerCmd(n),
		chat.NewChatCmd(n),
	)
}
