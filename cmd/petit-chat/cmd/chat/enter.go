package chat

import (
	"bufio"
	"fmt"

	"github.com/h0n9/petit-chat/util"
)

const (
	DefaultMsgTextEncoding = "UTF-8"
)

var enterCmd = util.NewCmd(
	"enter",
	"enter to chat",
	enterFunc,
)

func enterFunc(reader *bufio.Reader) error {
	// get user input
	fmt.Printf("Type chat room name: ")
	topic, err := util.GetInput(reader, false, false)
	if err != nil {
		return err
	}

	err = cli.StartChat(topic, reader)
	if err != nil {
		return err
	}

	return nil
}
