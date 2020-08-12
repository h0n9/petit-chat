package peer

import (
	"bufio"

	"github.com/h0n9/petit-chat/util"
)

var blackCmd = util.NewCmd(
	"black",
	"(unsupported) add a specific peer to black list",
	blackFunc,
)

func blackFunc(reader *bufio.Reader) error {
	// TODO: implement black peer

	return nil
}
