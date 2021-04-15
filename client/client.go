package client

import (
	"context"

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

func (c *Client) CreateMsgBox(tStr, nickname string) (*msg.Box, error) {
	topic, err := c.node.Join(tStr)
	if err != nil {
		return nil, err
	}
	// TODO: get metdata from parameters
	persona, err := types.NewPersona(nickname, []byte{}, c.node.PubKey)
	if err != nil {
		return nil, err
	}
	return c.msgCenter.CreateBox(topic, &persona)
}

func (c *Client) LeaveMsgBox(topicStr string) error {
	return c.msgCenter.LeaveBox(topicStr)
}

func (c *Client) GetMsgBox(topicStr string) (*msg.Box, bool) {
	return c.msgCenter.GetBox(topicStr)
}
