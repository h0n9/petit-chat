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
	msgCenter, err := node.GetCenter(hostPeer.GetNickname())
	if err != nil {
		return err
	}

	// get user input
	fmt.Printf("Type chat room name:")
	data, err := util.GetInput(reader)
	if err != nil {
		return err
	}

	_, err = msgCenter.CreateBox(data)
	if err != nil {
		return err
	}

	return nil
}
