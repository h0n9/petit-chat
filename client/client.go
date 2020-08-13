package client

import (
	"context"

	"github.com/h0n9/petit-chat/msg"
	"github.com/h0n9/petit-chat/p2p"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type Client struct {
	nickname string

	node      *p2p.Node
	msgCenter *msg.Center
	cfg       util.Config
}

func NewClient(ctx context.Context, cfg util.Config) (*Client, error) {
	node, err := p2p.NewNode(ctx, cfg)
	if err != nil {
		return nil, err
	}
	msgCenter, err := msg.NewCenter()
	if err != nil {
		return nil, err
	}
	return &Client{
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

func (c *Client) CreateMsgBox(topic string) error {
	return nil
}

func (c *Client) EnterMsgBox(topic string) error {
	return nil
}

func (c *Client) LeaveMsgBox(topic string) error {
	return nil
}
