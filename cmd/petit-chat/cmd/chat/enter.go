package chat

import (
	"bufio"
	"fmt"
	"strings"
	"sync"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/msg"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

const (
	DefaultMsgTextEncoding = "UTF-8"
)

var enterCmd = util.NewCmd(
	"enter",
	"enter to chat",
	enterFunc,
)

func enterFunc(reader *bufio.Reader) error {
	// get user input
	fmt.Printf("Type chat room name: ")
	topic, err := util.GetInput(reader, false, false)
	if err != nil {
		return err
	}

	msgBox, exist := cli.GetMsgBox(topic)
	if !exist {
		fmt.Printf("Type nickname: ")
		nickname, err := util.GetInput(reader, false, false)
		if err != nil {
			return err
		}
		fmt.Printf("Type public('true', 't' or 'false', 'f'): ")
		pubStr, err := util.GetInput(reader, false, true)
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
		err := readMsg(msgBox, msg)
		if err != nil {
			fmt.Printf("%s\n> ", err)
		}
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
			err  error    = nil
			stop bool     = false
			m    *msg.Msg = nil
		)
		for {
			select {
			case err = <-errs:
				fmt.Printf("%s\n> ", err)
			case m = <-msgSubCh:
				printMsg(msgBox, m)
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
			input, err := util.GetInput(reader, false, true)
			if err != nil {
				errs <- err
				continue
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
			case "/grant":
				fmt.Printf("<address> <R|W|X>: ")
				input, err := util.GetInput(reader, false, false)
				if err != nil {
					errs <- err
					continue
				}

				// parse strings
				strs := strings.Split(input, " ")
				if len(strs) != 2 {
					continue
				}
				addr := crypto.Addr(strs[0])
				if len(addr) != crypto.AddrSize {
					errs <- code.ImproperAddress
					continue
				}
				r, w, x := parsePerm(strs[1])

				err = msgBox.Grant(addr, r, w, x)
				if err != nil {
					errs <- err
					continue
				}
				continue
			case "/revoke":
				fmt.Printf("<address>: ")
				input, err := util.GetInput(reader, false, false)
				if err != nil {
					errs <- err
					continue
				}

				addr := crypto.Addr(input)
				if len(addr) != crypto.AddrSize {
					errs <- code.ImproperAddress
					continue
				}

				err = msgBox.Revoke(addr)
				if err != nil {
					errs <- err
					continue
				}
				continue
			case "":
				continue
			}
			if stop {
				break
			}

			// CLI supports ONLY TypeText
			msg := msg.NewMsgRaw(msgBox, types.EmptyHash, []byte(input), nil)
			err = msgBox.Publish(msg, true)
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
	if a.IsPublic() {
		p = "public"
	}
	str := fmt.Sprintf("Auth: %s\n", p)
	if len(a.Perms) > 0 {
		str += "Perms:\n"
	}
	for addr := range a.Perms {
		str += fmt.Sprintf("[%s] ", addr)
		if a.CanRead(addr) {
			str += "R"
		}
		if a.CanWrite(addr) {
			str += "W"
		}
		if a.CanExecute(addr) {
			str += "X"
		}
		str += "\n"
	}
	fmt.Printf("%s", str)
}

func printPeer(p *types.Persona) {
	fmt.Printf("[%s] %s\n", p.Address, p.Nickname)
}

func readMsg(b *msg.Box, m *msg.Msg) error {
	printMsg(b, m)
	if m.GetType() <= msg.TypeMeta {
		return nil
	}
	if m.GetClientAddr() == b.GetHostPersona().Address {
		return nil
	}
	meta := types.NewMeta(false, true, false)
	msgMeta := msg.NewMsgMeta(b, types.EmptyHash, m.GetHash(), meta)
	err := b.Publish(msgMeta, true)
	if err != nil {
		return err
	}
	return nil
}

func printMsg(box *msg.Box, m *msg.Msg) {
	timestamp := m.GetTimestamp()
	addr := m.GetSignature().PubKey.Address()
	persona := box.GetPersona(addr)
	nickname := "somebody"
	if persona != nil {
		nickname = persona.GetNickname()
	}
	switch m.GetType() {
	case msg.TypeRaw:
		body := m.GetBody().(msg.BodyRaw)
		metas := m.GetMetas()
		fmt.Printf("[%s, %s] %s\n", timestamp, nickname, body.Data)
		for addr, meta := range metas {
			nickname = box.GetPersona(addr).Nickname
			fmt.Printf("  - %s %s\n", nickname, printMeta(meta))
		}
	case msg.TypeHelloSyn:
		fmt.Printf("[%s, %s] entered\n", timestamp, nickname)
	case msg.TypeHelloAck:
	case msg.TypeBye:
		fmt.Printf("[%s, %s] left\n", timestamp, nickname)
	case msg.TypeUpdate:
	case msg.TypeMeta:
		// body := m.GetBody().(msg.BodyMeta)
		// done := printMeta(body.Meta)
		// fmt.Printf("[%s, %s] %s %x\n", timestamp, nickname, done, m.GetParentHash())
	default:
		fmt.Println("Unknown Type")
	}
}

func printMeta(meta types.Meta) string {
	str := ""
	if meta.Received() {
		str += "received,"
	}
	if meta.Read() {
		str += "read,"
	}
	if meta.Typing() {
		str += "typing,"
	}
	return str
}

func parsePerm(permStr string) (bool, bool, bool) {
	r, w, x := false, false, false
	permStr = strings.ToUpper(permStr)
	if strings.Contains(permStr, "R") {
		r = true
	}
	if strings.Contains(permStr, "W") {
		w = true
	}
	if strings.Contains(permStr, "X") {
		x = true
	}
	return r, w, x
}
