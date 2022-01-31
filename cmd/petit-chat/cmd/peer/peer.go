package peer

import (
	"github.com/h0n9/petit-chat/client"
	"github.com/h0n9/petit-chat/server"
	"github.com/h0n9/petit-chat/util"
)

var (
	svr *server.Server
	cli *client.Client
)

func NewPeerCmd(s *server.Server, c *client.Client) *util.Cmd {
	svr = s
	cli = c
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
