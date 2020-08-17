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
	data, err := util.GetInput(reader, false)
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

	errs := make(chan error, 1)
	go func() {
		for {
			fmt.Printf("> ")
			data, err = util.GetInput(reader, false)
			if err != nil {
				errs <- err
				return
			}
			switch data {
			case "/exit":
				return
			case "/close":
				err = msgBox.Close()
				if err != nil {
					errs <- err
				}
				return
			case "/msgs":
				msgs := msgBox.GetMsgs()
				for time, msg := range msgs {
					fmt.Println("time:", time, ", from:", msg.GetFrom(), ",", msg.GetData())
				}
				continue
			case "":
				continue
			}

			err = msgBox.Publish([]byte(data))
			if err != nil {
				errs <- err
				return
			}
		}
	}()

	err = msgBox.Open()
	if err != nil {
		return err
	}

	return nil
}
