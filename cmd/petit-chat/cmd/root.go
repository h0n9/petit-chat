package cmd

import (
	"github.com/h0n9/petit-chat/client"
	"github.com/h0n9/petit-chat/cmd/petit-chat/cmd/chat"
	"github.com/h0n9/petit-chat/cmd/petit-chat/cmd/peer"
	"github.com/h0n9/petit-chat/server"
	"github.com/h0n9/petit-chat/util"
)

var (
	svr *server.Server
	cli *client.Client
)

func NewRootCmd(s *server.Server, c *client.Client) *util.Cmd {
	svr = s
	cli = c
	return util.NewCmd("petit-chat", "entry point for petit-chat", nil,
		infoCmd,
		peer.NewPeerCmd(s, c),
		chat.NewChatCmd(s, c),
	)
}
