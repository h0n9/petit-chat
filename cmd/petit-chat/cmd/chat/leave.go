package chat

import (
	"bufio"
	"fmt"

	"github.com/h0n9/petit-chat/util"
)

var leaveCmd = util.NewCmd(
	"leave",
	"leave chat",
	leaveFunc,
)

func leaveFunc(reader *bufio.Reader) error {
	// get user input
	fmt.Printf("Type chat room name: ")
	data, err := util.GetInput(reader, false)
	if err != nil {
		return err
	}

	err = cli.LeaveMsgBox(data)
	if err != nil {
		return err
	}

	return nil
}
