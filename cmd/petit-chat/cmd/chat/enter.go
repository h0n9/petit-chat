package chat

import (
	"bufio"
	"fmt"
	"sync"

	"github.com/h0n9/petit-chat/msg"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

const (
	DEFAULT_MSG_TEXT_ENCODING = "UTF-8"
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
		fmt.Printf("Type public('true', 't' or 'false', 'f'): ")
		pubStr, err := util.GetInput(reader, false)
		if err != nil {
			return err
		}
		pub, err := util.ToBool(pubStr)
		if err != nil {
			return err
		}
		mb, err := cli.CreateMsgBox(topic, nickname, pub)
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
			input, err := util.GetInput(reader, false)
			if err != nil {
				errs <- err
				return
			}
			switch input {
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
			case "/auth":
				auth := msgBox.GetAuth()
				printAuth(auth)
				continue
			case "":
				continue
			}
			if stop {
				break
			}

			// encapulate user input into MsgStructText
			mst := msg.NewMsgStructText([]byte(input), DEFAULT_MSG_TEXT_ENCODING)
			data, err := mst.Encapsulate()
			if err != nil {
				errs <- err
			}

			// CLI supports ONLY MsgTypeText
			err = msgBox.Publish(msg.MsgTypeText, types.Hash{}, true, data)
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

func printAuth(a *types.Auth) {
	p := "private"
	if a.IsPublic {
		p = "public"
	}
	str := fmt.Sprintf("Auth: %s\n", p)
	if len(a.Perms) > 0 {
		str += "Perms:\n"
	}
	for id, perm := range a.Perms {
		str += fmt.Sprintf("[%s] ", id)
		if perm.Read {
			str += "R"
		}
		if perm.Write {
			str += "W"
		}
		if perm.Execute {
			str += "X"
		}
		str += "\n"
	}
	fmt.Printf("%s", str)
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
	case msg.MsgTypeText:
		mst := msg.NewMsgStructText(nil, DEFAULT_MSG_TEXT_ENCODING)
		err := mst.Decapsulate(m.GetData())
		if err != nil {
			return
		}
		fmt.Printf("[%s, %s] %s\n", timestamp, nickname, mst.GetData())
	case msg.MsgTypeImage:
		// TODO: CLI doesn't support this type
	case msg.MsgTypeVideo:
		// TODO: CLI doesn't support this type
	case msg.MsgTypeAudio:
		// TODO: CLI doesn't support this type
	case msg.MsgTypeRaw:
		// TODO: CLI doesn't support this type
	case msg.MsgTypeHello:
		if m.ParentMsgHash.IsEmpty() {
			fmt.Printf("[%s, %s] entered\n", timestamp, nickname)
		}
	case msg.MsgTypeBye:
		// do nothing
	default:
		fmt.Println("Unknown MsgType")
	}
}
