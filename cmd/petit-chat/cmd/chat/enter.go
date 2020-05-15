package chat

import (
	"bufio"

	"github.com/h0n9/petit-chat/util"
)

var enterCmd = util.NewCmd(
	"enter",
	"enter to chat",
	enterFunc,
)

func enterFunc(reader *bufio.Reader) error {

	return nil
}
