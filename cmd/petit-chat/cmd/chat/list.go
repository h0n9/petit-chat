package chat

import (
	"bufio"
	"fmt"

	"github.com/h0n9/petit-chat/client"
	"github.com/h0n9/petit-chat/util"
)

var listCmd = util.NewCmd(
	"list",
	"show a list of chats",
	listFunc,
)

func listFunc(reader *bufio.Reader) error {
	chats := cli.GetChats()
	printChats(chats)
	return nil
}

func printChats(chats map[string]*client.Chat) {
	if len(chats) == 0 {
		fmt.Printf("none\n")
		return
	}
	n := 1
	for topic := range chats {
		fmt.Printf("%d. %s\n", n, topic)
		n++
	}
}
