package chat

import (
	"bufio"
	"fmt"
	"sync"

	"github.com/h0n9/petit-chat/msg"
	"github.com/h0n9/petit-chat/types"
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
	topic, err := util.GetInput(reader, false)
	if err != nil {
		return err
	}

	msgBox, exist := cli.GetMsgBox(topic)
	if !exist {
		fmt.Printf("Type nickname: ")
		nickname, err := util.GetInput(reader, false)
		if err != nil {
			return err
		}
		mb, err := cli.CreateMsgBox(topic, nickname)
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
		printMsg(msgBox, msg)
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
				printMsg(msgBox, msg)
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
			data, err := util.GetInput(reader, false)
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
					printMsg(msgBox, msg)
				}
				continue
			case "/peers":
				peers := msgBox.GetPersonae()
				for _, peer := range peers {
					printPeer(peer)
				}
				continue
			case "":
				continue
			}
			if stop {
				break
			}

			// CLI supports ONLY MsgTypeText
			err = msgBox.Publish(types.MsgText, types.Hash{}, []byte(data))
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

func printPeer(p *types.Persona) {
	fmt.Printf("[%s] %s\n", p.Address, p.Nickname)
}

func printMsg(b *msg.Box, m *msg.Msg) {
	timestamp := m.GetTime()
	from := m.GetFrom()
	persona := b.GetPersona(from.ClientAddr)
	nickname := "somebody"
	if persona != nil {
		nickname = persona.GetNickname()
	}
	switch m.GetType() {
	case types.MsgText:
		fmt.Printf("[%s, %s] %s\n", timestamp, nickname, string(m.GetData()))
	case types.MsgImage:
		// TODO: CLI doesn't support this type
	case types.MsgVideo:
		// TODO: CLI doesn't support this type
	case types.MsgAudio:
		// TODO: CLI doesn't support this type
	case types.MsgRaw:
		// TODO: CLI doesn't support this type
	case types.MsgHello:
		if types.IsEmpty(m.ParentMsgHash) {
			fmt.Printf("[%s, %s] entered\n", timestamp, nickname)
		}
	case types.MsgBye:
		// do nothing
	default:
		fmt.Println("Unknown MsgType")
	}
}
