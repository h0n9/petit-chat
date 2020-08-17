package chat

import (
	"bufio"
	"fmt"
	"sync"

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

	var (
		wait sync.WaitGroup
		stop bool = false
	)
	wait.Add(1)

	// TODO: Fix error handling of goroutines
	errs := make(chan error, 1)
	defer close(errs)

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
				stop = true
			case "/msgs":
				msgs := msgBox.GetMsgs()
				for time, msg := range msgs {
					fmt.Printf("[%s, %s] %s\n", time, msg.GetFrom(), string(msg.GetData()))
				}
				continue
			case "":
				continue
			}

			if stop {
				break
			}

			err = msgBox.Publish([]byte(data))
			if err != nil {
				errs <- err
				return
			}
		}
		defer wait.Done()
	}()

	go msgBox.Subscribe()

	wait.Wait()

	return nil
}
