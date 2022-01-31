package cmd

import (
	"bufio"

	"github.com/h0n9/petit-chat/util"
)

var infoCmd = util.NewCmd(
	"info",
	"current node's brief information",
	infoFunc,
)

func infoFunc(reader *bufio.Reader) error {
	svr.PrintInfo()
	return nil
}
