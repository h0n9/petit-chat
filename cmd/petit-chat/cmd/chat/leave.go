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

	err = msgCenter.LeaveBox(data)
	if err != nil {
		return err
	}

	return nil
}
