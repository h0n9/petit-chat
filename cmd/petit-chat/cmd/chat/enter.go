package chat

import (
	"bufio"
	"fmt"

	"github.com/h0n9/petit-chat/msg"
	"github.com/h0n9/petit-chat/util"
)

var enterCmd = util.NewCmd(
	"enter",
	"enter to chat",
	enterFunc,
)

func enterFunc(reader *bufio.Reader) error {
	Center, err := node.GetCenter(hostPeer.GetNickname())
	if err != nil {
		return err
	}

	// get user input
	fmt.Printf("Type chat room name:")
	data, err := util.GetInput(reader)
	if err != nil {
		return err
	}

	msgBox, err := Center.GetBox(data)
	if err != nil {
		return err
	}

	interact(msgBox)

	return nil
}

func interact(msgBox *msg.Box) {
	// interact with msgBox
	// expected features:
	// - send msg
	// - receive msg
	// - etc
}
