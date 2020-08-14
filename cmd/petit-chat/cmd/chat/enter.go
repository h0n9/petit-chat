package chat

import (
	"bufio"
	"fmt"

	"github.com/h0n9/petit-chat/util"
)

var enterCmd = util.NewCmd(
	"enter",
	"enter to chat",
	enterFunc,
)

func enterFunc(reader *bufio.Reader) error {
	// get user input
	fmt.Printf("Type chat room name: ")
	data, err := util.GetInput(reader, true)
	if err != nil {
		return err
	}

	msgBox, exist := cli.GetMsgBox(data)
	if !exist {
		mb, err := cli.CreateMsgBox(data)
		if err != nil {
			return err
		}
		msgBox = mb
	}

	sub, err := msgBox.Subscribe()
	if err != nil {
		return err
	}
	defer sub.Cancel()

	ctx := cli.GetContext()
	myID := cli.GetID()
	go func() {
		for {
			received, err := sub.Next(ctx)
			if err != nil {
				return
			}
			if received.GetFrom() != myID {
				fmt.Printf("%s> %s\n", received.GetFrom(), received.GetData())
			}
			fmt.Printf("> ")
		}
	}()

	fmt.Printf("> ")
	for {
		data, err = util.GetInput(reader, false)
		if err != nil {
			return err
		}
		if data == "/exit" {
			break
		}
		if data == "" {
			fmt.Printf("> ")
			continue
		}

		msgBox.Publish([]byte(data))
	}

	return nil
}
