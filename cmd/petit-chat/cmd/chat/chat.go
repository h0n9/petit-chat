package chat

import (
	"github.com/h0n9/petit-chat/client"
	"github.com/h0n9/petit-chat/util"
)

var (
	cli *client.Client
)

func NewChatCmd(c *client.Client) *util.Cmd {
	cli = c
	return util.NewCmd(
		"chat",
		"chat related commands",
		nil,
		listCmd,
		enterCmd,
		leaveCmd,
	)
}
