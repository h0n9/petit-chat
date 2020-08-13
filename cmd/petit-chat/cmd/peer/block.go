package peer

import (
	"bufio"

	"github.com/h0n9/petit-chat/util"
)

var blockCmd = util.NewCmd(
	"block",
	"(unsupported) add a specific peer to block list",
	blockFunc,
)

func blockFunc(reader *bufio.Reader) error {
	// TODO: implement block peer

	return nil
}
