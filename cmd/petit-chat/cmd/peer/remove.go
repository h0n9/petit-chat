package peer

import (
	"bufio"

	"github.com/h0n9/petit-chat/util"
)

var removeCmd = util.NewCmd(
	"remove",
	"(unsupported) remove a specific peer from peer list",
	removeFunc,
)

func removeFunc(reader *bufio.Reader) error {
	// TODO: implement remove peer

	return nil
}
