package chat

import (
	"bufio"
	"fmt"

	"github.com/h0n9/petit-chat/util"
)

var listCmd = util.NewCmd(
	"list",
	"show a list of chats",
	listFunc,
)

func listFunc(reader *bufio.Reader) error {
	msgCenter := node.GetMsgCenter()

	for topic, msgBox := range msgCenter.GetMsgBoxes() {
		fmt.Printf("%s\n", topic)
		for _, p := range msgBox.GetPeers() {
			fmt.Printf(" - %s\n", p)
		}
	}

	return nil
}
