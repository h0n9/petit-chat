package chat

import (
	"bufio"
	"fmt"
	"sync"

	"github.com/h0n9/petit-chat/msg"
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

	// TODO: Fix error handling of goroutines
	errs := make(chan error, 1)
	defer close(errs)

	// open subscription
	go msgBox.Subscribe(msg.DefaultMsgHandler)

	// get and print out received msgs
	msgs := msgBox.GetUnreadMsgs()
	for _, msg := range msgs {
		printMsg(msg)
	}

	// get and print out new msgs
	var (
		msgSubCh     = make(chan *msg.Msg, 1)
		msgStopSubCh = make(chan bool, 1)
	)
	defer close(msgSubCh)
	defer close(msgStopSubCh)

	msgBox.SetMsgSubCh(msgSubCh)
	defer msgBox.SetMsgSubCh(nil)

	wait.Add(1)
	go func() {
		var (
			stop bool     = false
			msg  *msg.Msg = nil
		)
		for {
			select {
			case msg = <-msgSubCh:
				printMsg(msg)
			case <-msgStopSubCh:
				stop = true
			}
			if stop {
				break
			}
		}
		wait.Done()
	}()

	// get user input
	wait.Add(1)
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
				msgStopSubCh <- true
				stop = true
			case "/msgs":
				msgs := msgBox.GetMsgs()
				for _, msg := range msgs {
					printMsg(msg)
				}
				continue
			case "":
				continue
			}
			if stop {
				break
			}

			// CLI supports ONLY MsgTypeText
			err = msgBox.Publish(msg.MsgTypeText, []byte(data))
			if err != nil {
				errs <- err
				return
			}
		}
		wait.Done()
	}()

	wait.Wait()

	return nil
}

func printMsg(m *msg.Msg) {
	switch m.GetType() {
	case msg.MsgTypeText:
		fmt.Printf("[%s, %s] %s\n", m.GetTime(), m.GetFrom(), string(m.GetData()))
	case msg.MsgTypeImage:
		// TODO: CLI doesn't support this type
	case msg.MsgTypeVideo:
		// TODO: CLI doesn't support this type
	case msg.MsgTypeAudio:
		// TODO: CLI doesn't support this type
	case msg.MsgTypeRaw:
		// TODO: CLI doesn't support this type
	default:
		fmt.Println("Unknowm MsgType")
	}
}
