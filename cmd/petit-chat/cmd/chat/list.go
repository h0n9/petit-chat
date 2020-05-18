package chat

import (
	"bufio"
	"fmt"
	"strconv"

	"github.com/h0n9/petit-chat/msg"
	"github.com/h0n9/petit-chat/util"
)

var listCmd = util.NewCmd(
	"list",
	"show a list of chats",
	listFunc,
)

func listFunc(reader *bufio.Reader) error {
	printMsgBoxes()
	return nil
}

func printMsgBoxes() map[string]*msg.MsgBox {
	msgBoxes := node.GetMsgCenter().GetMsgBoxes()
	n := 1
	for topic, msgBox := range msgBoxes {
		nStr := strconv.Itoa(n)
		n += 1
		msgBoxes[nStr] = msgBox

		fmt.Printf("%s\n", topic)
		for _, p := range msgBox.GetPeers() {
			fmt.Printf(" - %s\n", p)
		}
	}
	return msgBoxes
}
