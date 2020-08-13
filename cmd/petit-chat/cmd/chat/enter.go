package chat

import (
	"bufio"
	"fmt"

	"github.com/h0n9/petit-chat/util"
)

var enterCmd = util.NewCmd(
	"enter",
	"enter to chat",
	enterFunc,
)

func enterFunc(reader *bufio.Reader) error {
	// get user input
	fmt.Printf("Type chat room name:")
	data, err := util.GetInput(reader)
	if err != nil {
		return err
	}

	err = cli.EnterMsgBox(data)
	if err != nil {
		return err
	}

	return nil
}
