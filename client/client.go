package client

import (
	"bufio"
	"fmt"

	"github.com/h0n9/petit-chat/server"
	"github.com/h0n9/petit-chat/util"
)

type Client struct {
	svr *server.Server
}

func NewClient(svr *server.Server) (*Client, error) {
	return &Client{
		svr: svr,
	}, nil
}

func (c *Client) StartChat(topic string, reader *bufio.Reader) error {
	box, exist := c.svr.GetMsgBox(topic)
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
		b, err := c.svr.CreateMsgBox(topic, nickname, pub)
		if err != nil {
			return err
		}
		box = b
	}

	chat, err := NewChat(box, reader)
	if err != nil {
		return err
	}
	defer chat.Close()

	// open subscription
	go box.Subscribe()

	// get and print out received msgs
	msgs := box.GetUnreadMsgs()
	for _, msg := range msgs {
		err := readMsg(box, msg)
		if err != nil {
			fmt.Printf("%s\n> ", err)
		}
	}

	// get and print out new msgs
	chat.wg.Add(1)
	go chat.Receive()

	// get user input
	chat.wg.Add(1)
	go chat.Send()

	chat.wg.Wait()

	return nil
}
