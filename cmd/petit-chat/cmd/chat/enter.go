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
	msgBoxes := node.GetMsgCenter().GetMsgBoxes()
	printMsgBoxes(msgBoxes)
	if len(msgBoxes) == 0 {
		return nil
	}

	// get user input
	data, err := util.GetInput(reader)
	if err != nil {
		return err
	}

	msgBox, exist := msgBoxes[data]
	if !exist {
		return fmt.Errorf("'%s' not proper chat room", data)
	}

	interact(msgBox)

	return nil
}

func interact(msgBox *msg.MsgBox) {
	// interact with msgBox
	// expected features:
	// - send msg
	// - receive msg
	// - etc
}
