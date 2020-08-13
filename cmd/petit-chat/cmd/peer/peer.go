package peer

import (
	"github.com/h0n9/petit-chat/client"
	"github.com/h0n9/petit-chat/util"
)

var (
	cli *client.Client
)

func NewPeerCmd(c *client.Client) *util.Cmd {
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
