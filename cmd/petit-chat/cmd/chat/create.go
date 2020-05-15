package chat

import (
	"bufio"

	"github.com/h0n9/petit-chat/util"
)

var createCmd = util.NewCmd(
	"create",
	"create to chat",
	createFunc,
)

func createFunc(reader *bufio.Reader) error {

	return nil
}
