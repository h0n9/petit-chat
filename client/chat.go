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
	box   *msg.Box
	vault *types.Vault

	wg              sync.WaitGroup
	chStopReceive   chan bool
	chError         chan error
	chMsgCapsuleSub chan *msg.MsgCapsule

	reader *bufio.Reader
}

func NewChat(box *msg.Box, vault *types.Vault, reader *bufio.Reader) (*Chat, error) {
	return &Chat{
		box:   box,
		vault: vault,

		wg:              sync.WaitGroup{},
		chStopReceive:   make(chan bool, 1),
		chError:         make(chan error, 1),
		chMsgCapsuleSub: nil,

		reader: reader,
	}, nil
}

func (c *Chat) setChMsgCapsule(chMsgCapsule chan *msg.MsgCapsule) {
	c.chMsgCapsuleSub = chMsgCapsule
}

func (c *Chat) Close() {
	close(c.chStopReceive)
	close(c.chError)
}

func (c *Chat) Subscribe() error {
	c.setChMsgCapsule(c.box.GetChMsgCapsule())
	return c.box.Subscribe()
}

func (c *Chat) Stop() {
	c.setChMsgCapsule(nil)
}

func (c *Chat) Send() {
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
		peerID := c.box.GetHostID()
		clientAddr := c.vault.GetAddr()
		m := msg.NewMsgRaw(peerID, clientAddr, types.EmptyHash, []byte(input), nil)
		err = c.publish(m, false)
		if err != nil {
			c.chError <- err
			continue
		}
	}
	c.wg.Done()
}

func (c *Chat) Receive() {
	var (
		stop       bool  = false
		err        error = nil
		msgCapsule *msg.MsgCapsule
	)
	for {
		select {
		case msgCapsule = <-c.chMsgCapsuleSub:
			err := msgCapsule.Check()
			if err != nil {
				c.chError <- err
				continue
			}
			m, err := msgCapsule.Decapsulate(c.vault.GetSecretKey())
			if err != nil {
				c.chError <- err
				continue
			}
			err = m.Verify()
			if err != nil {
				c.chError <- err
				continue
			}
			printMsg(c.box, m)
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

func (c *Chat) publish(m *msg.Msg, encrypt bool) error {
	err := m.Sign(c.vault.GetPrivKey())
	if err != nil {
		return err
	}

	msgCapsule, err := m.Encapsulate(encrypt, c.vault.GetSecretKey())
	if err != nil {
		return err
	}

	err = c.box.Publish(msgCapsule)
	if err != nil {
		return err
	}

	return nil
}
