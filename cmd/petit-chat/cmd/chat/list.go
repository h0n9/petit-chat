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
	msgBoxes := node.GetMsgCenter().GetMsgBoxes()
	printMsgBoxes(msgBoxes)
	return nil
}

func printMsgBoxes(msgBoxes map[string]*msg.MsgBox) {
	n := 1
	for topic, msgBox := range msgBoxes {
		nStr := strconv.Itoa(n)
		msgBoxes[nStr] = msgBox

		fmt.Printf("%d. %s\n", n, topic)
		for _, p := range msgBox.GetPeers() {
			fmt.Printf(" - %s\n", p)
		}
		n++
	}
}
