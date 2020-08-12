package chat

import (
	"bufio"
	"fmt"

	"github.com/h0n9/petit-chat/msg"
	"github.com/h0n9/petit-chat/util"
)

var listCmd = util.NewCmd(
	"list",
	"show a list of chats",
	listFunc,
)

func listFunc(reader *bufio.Reader) error {
	msgCenter, err := node.GetCenter(hostPeer.GetNickname())
	if err != nil {
		return err
	}
	msgBoxes := msgCenter.GetBoxes()
	printBoxes(msgBoxes)
	return nil
}

func printBoxes(msgBoxes map[string]*msg.Box) {
	if len(msgBoxes) == 0 {
		fmt.Printf("none\n")
		return
	}
	n := 1
	for topic, msgBox := range msgBoxes {
		fmt.Printf("%d. %s\n", n, topic)
		for _, p := range msgBox.GetPeers() {
			fmt.Printf(" - %s\n", p.GetNickname())
		}
		n++
	}
}
