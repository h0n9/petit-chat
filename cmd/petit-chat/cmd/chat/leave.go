package chat

import (
	"bufio"

	"github.com/h0n9/petit-chat/util"
)

var leaveCmd = util.NewCmd(
	"leave",
	"leave chat",
	leaveFunc,
)

func leaveFunc(reader *bufio.Reader) error {

	return nil
}
