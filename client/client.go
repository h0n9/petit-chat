package client

import (
	"bufio"
	"context"
	"fmt"

	"github.com/h0n9/petit-chat/msg"
	"github.com/h0n9/petit-chat/p2p"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type Client struct {
	ctx context.Context

	nickname string

	node      *p2p.Node
	msgCenter *msg.Center
	cfg       *util.Config
}

func NewClient(ctx context.Context, cfg *util.Config) (*Client, error) {
	node, err := p2p.NewNode(ctx, cfg)
	if err != nil {
		return nil, err
	}
	msgCenter, err := msg.NewCenter(ctx, node.GetHostID())
	if err != nil {
		return nil, err
	}
	return &Client{
		ctx:       ctx,
		nickname:  "",
		node:      node,
		msgCenter: msgCenter,
		cfg:       cfg,
	}, nil
}

func (c *Client) Close() error {
	return c.node.Close()
}

func (c *Client) Info() {
	c.node.Info()
}

func (c *Client) GetID() types.ID {
	return c.node.GetHostID()
}

func (c *Client) GetNickname() string {
	return c.nickname
}

func (c *Client) DiscoverPeers() error {
	return c.node.DiscoverPeers(c.cfg.BootstrapNodes)
}

func (c *Client) GetPeers() []types.ID {
	return c.node.GetPeers()
}

func (c *Client) GetMsgCenter() *msg.Center {
	return c.msgCenter
}

func (c *Client) CreateMsgBox(tStr, nickname string, pub bool) (*msg.Box, error) {
	topic, err := c.node.Join(tStr)
	if err != nil {
		return nil, err
	}
	// TODO: get metdata from parameters
	persona, err := types.NewPersona(nickname, []byte{}, c.node.PubKey)
	if err != nil {
		return nil, err
	}
	return c.msgCenter.CreateBox(topic, pub, &c.node.PrivKey, &persona)
}

func (c *Client) LeaveMsgBox(topicStr string) error {
	return c.msgCenter.LeaveBox(topicStr)
}

func (c *Client) GetMsgBox(topicStr string) (*msg.Box, bool) {
	return c.msgCenter.GetBox(topicStr)
}

func (cli *Client) StartChat(topic string, reader *bufio.Reader) error {
	box, exist := cli.GetMsgBox(topic)
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
		b, err := cli.CreateMsgBox(topic, nickname, pub)
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
	go chat.Output()

	// get user input
	chat.wg.Add(1)
	go chat.Input()

	chat.wg.Wait()

	return nil
}
