package cmd

import (
	"github.com/h0n9/petit-chat/cmd/petit-chat/cmd/chat"
	"github.com/h0n9/petit-chat/cmd/petit-chat/cmd/peer"
	"github.com/h0n9/petit-chat/p2p"
	"github.com/h0n9/petit-chat/util"
)

func NewRootCmd(node *p2p.Node) *util.Cmd {
	return util.NewCmd("petit-chat", "entry point for petit-chat", nil,
		peer.NewPeerCmd(node),
		chat.NewChatCmd(node),
	)
}
