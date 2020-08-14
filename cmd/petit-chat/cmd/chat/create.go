package chat

import (
	"bufio"
	"fmt"

	"github.com/h0n9/petit-chat/util"
)

var createCmd = util.NewCmd(
	"create",
	"create to chat",
	createFunc,
)

func createFunc(reader *bufio.Reader) error {
	// get user input
	fmt.Printf("Type chat room name:")
	data, err := util.GetInput(reader, true)
	if err != nil {
		return err
	}

	err = cli.CreateMsgBox(data)
	if err != nil {
		return err
	}

	return nil
}
