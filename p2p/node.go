package p2p

import (
	"context"
	"fmt"

	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/util"
)

type Node struct {
	ctx context.Context

	PrivKey crypto.PrivKey
	PubKey  crypto.PubKey
	Address crypto.Addr

	Host   Host
	PubSub PubSub
}

func NewNode(cfg util.Config) (Node, error) {
	node := Node{}

	node.ctx = context.Background()

	privKey, err := crypto.GenPrivKey()
	if err != nil {
		return Node{}, nil
	}

	node.PrivKey = privKey
	node.PubKey = privKey.PubKey()
	node.Address = node.PubKey.Address()

	host, err := NewHost(node.ctx, node.PrivKey, cfg.ListenAddrs)
	if err != nil {
		return Node{}, nil
	}

	node.Host = host

	return node, nil
}

func (n *Node) Info() {
	if n.Host == nil {
		return
	}

	fmt.Println("address:", n.Address)
	fmt.Println("host ID:", n.Host.ID().Pretty())
	fmt.Println("host addrs:", n.Host.Addrs())

	fmt.Printf("%s/p2p/%s\n", n.Host.Addrs()[0], n.Host.ID())
}
