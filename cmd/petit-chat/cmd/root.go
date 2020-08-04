package cmd

import (
	"github.com/h0n9/petit-chat/cmd/petit-chat/cmd/chat"
	"github.com/h0n9/petit-chat/cmd/petit-chat/cmd/peer"
	"github.com/h0n9/petit-chat/msg"
	"github.com/h0n9/petit-chat/p2p"
	"github.com/h0n9/petit-chat/util"
)

var (
	node     *p2p.Node
	hostPeer msg.Peer
)

func NewRootCmd(n *p2p.Node, h msg.Peer) *util.Cmd {
	node = n
	hostPeer = h
	return util.NewCmd("petit-chat", "entry point for petit-chat", nil,
		infoCmd,
		peer.NewPeerCmd(n, h),
		chat.NewChatCmd(n, h),
	)
}
