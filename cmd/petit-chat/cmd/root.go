package cmd

import (
	"github.com/h0n9/petit-chat/client"
	"github.com/h0n9/petit-chat/cmd/petit-chat/cmd/chat"
	"github.com/h0n9/petit-chat/cmd/petit-chat/cmd/peer"
	"github.com/h0n9/petit-chat/util"
)

var (
	cli *client.Client
)

func NewRootCmd(c *client.Client) *util.Cmd {
	cli = c
	return util.NewCmd("petit-chat", "entry point for petit-chat", nil,
		infoCmd,
		peer.NewPeerCmd(c),
		chat.NewChatCmd(c),
	)
}
