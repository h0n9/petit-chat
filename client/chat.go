package client

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

type Chat struct {
	box *msg.Box

	wg              sync.WaitGroup
	chStopReceive   chan bool
	chError         chan error
	chMsgCapsuleSub chan msg.MsgCapsule

	reader *bufio.Reader
}

func NewChat(box *msg.Box, reader *bufio.Reader) (*Chat, error) {
	return &Chat{
		box: box,

		wg:              sync.WaitGroup{},
		chStopReceive:   make(chan bool, 1),
		chError:         make(chan error, 1),
		chMsgCapsuleSub: box.GetChMsgCapsule(),

		reader: reader,
	}, nil
}

func (c *Chat) SetChMsgCapsule(chMsgCapsule chan msg.MsgCapsule) {
	c.chMsgCapsuleSub = chMsgCapsule
}

func (c *Chat) Close() {
	close(c.chStopReceive)
	close(c.chError)
}

func (c *Chat) input() {
	var stop bool = false
	for {
		fmt.Printf("> ")
		input, err := util.GetInput(c.reader, false, true)
		if err != nil {
			c.chError <- err
			continue
		}
		switch input {
		case "/exit":
			c.chStopReceive <- true
			stop = true
		case "/msgs":
			msgs := c.box.GetMsgs()
			for _, msg := range msgs {
				printMsg(c.box, msg)
			}
			continue
		case "/peers":
			peers := c.box.GetPersonae()
			for _, peer := range peers {
				printPeer(peer)
			}
			continue
		case "/auth":
			auth := c.box.GetAuth()
			printAuth(auth)
			continue
		case "/grant":
			fmt.Printf("<address> <R|W|X>: ")
			input, err := util.GetInput(c.reader, false, false)
			if err != nil {
				c.chError <- err
				continue
			}

			// parse strings
			strs := strings.Split(input, " ")
			if len(strs) != 2 {
				continue
			}
			addr := crypto.Addr(strs[0])
			if len(addr) != crypto.AddrSize {
				c.chError <- code.ImproperAddress
				continue
			}
			r, w, x := parsePerm(strs[1])

			err = c.box.Grant(addr, r, w, x)
			if err != nil {
				c.chError <- err
				continue
			}
			continue
		case "/revoke":
			fmt.Printf("<address>: ")
			input, err := util.GetInput(c.reader, false, false)
			if err != nil {
				c.chError <- err
				continue
			}

			addr := crypto.Addr(input)
			if len(addr) != crypto.AddrSize {
				c.chError <- code.ImproperAddress
				continue
			}

			err = c.box.Revoke(addr)
			if err != nil {
				c.chError <- err
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
		msg := msg.NewMsgRaw(c.box, types.EmptyHash, []byte(input), nil)
		err = c.box.Publish(msg, true)
		if err != nil {
			c.chError <- err
			return
		}
	}
	c.wg.Done()
}

func (c *Chat) output() {
	var (
		stop       bool  = false
		err        error = nil
		msgCapsule msg.MsgCapsule
	)
	for {
		select {
		case msgCapsule = <-c.chMsgCapsuleSub:
			fmt.Printf("%s\n", string(msgCapsule.Data))
			// TODO: handler comes here
		case err = <-c.chError:
			fmt.Println(err)
		case <-c.chStopReceive:
			stop = true
		}
		if stop {
			break
		}
	}
	c.wg.Done()
}
