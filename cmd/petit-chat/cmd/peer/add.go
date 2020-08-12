package peer

import (
	"bufio"

	"github.com/h0n9/petit-chat/util"
)

var addCmd = util.NewCmd(
	"add",
	"(unsupported) add a specific peer to peer list",
	addFunc,
)

func addFunc(reader *bufio.Reader) error {
	// TODO: implement add peer

	return nil
}
