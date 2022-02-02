package client

import (
	"bufio"
	"fmt"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/server"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type Client struct {
	svr   *server.Server
	chats map[string]*Chat
}

func NewClient(svr *server.Server) (*Client, error) {
	return &Client{
		svr:   svr,
		chats: make(map[string]*Chat),
	}, nil
}

func (c *Client) GetChats() map[string]*Chat {
	return c.chats
}

func (c *Client) GetChat(topic string) (*Chat, bool) {
	chat, exist := c.chats[topic]
	return chat, exist
}

func (c *Client) SetChat(topic string, chat *Chat) error {
	_, exist := c.GetChat(topic)
	if exist {
		return code.AlreadyExistingTopic
	}
	c.chats[topic] = chat
	return nil
}

func (c *Client) StartChat(topic string, reader *bufio.Reader) error {
	chat, exist := c.GetChat(topic)
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
		// init vault
		privKey, err := crypto.GenPrivKey()
		if err != nil {
			return err
		}
		secretKey, err := crypto.GenSecretKey()
		if err != nil {
			return err
		}
		persona, err := types.NewPersona(nickname, nil, privKey.PubKey())
		if err != nil {
			return err
		}
		vault := types.NewVault(persona, privKey, secretKey)
		box, err := c.svr.CreateMsgBox(topic, nickname, pub)
		if err != nil {
			return err
		}
		newChat, err := NewChat(box, vault, reader)
		if err != nil {
			return err
		}
		err = c.SetChat(topic, newChat)
		if err != nil {
			return err
		}
		chat = newChat
	}

	// open subscription
	go chat.Subscribe()
	defer chat.Stop()

	// start goroutine for receiving msgs
	chat.wg.Add(1)
	go chat.Receive()

	// start goroutine for sending msgs
	chat.wg.Add(1)
	go chat.Send()

	// wait for all of goroutines to stop
	chat.wg.Wait()

	return nil
}
