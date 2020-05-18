package p2p

import (
	"context"
	"fmt"

	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/msg"
	"github.com/h0n9/petit-chat/util"
)

type Node struct {
	ctx context.Context

	PrivKey crypto.PrivKey
	PubKey  crypto.PubKey
	Address crypto.Addr

	host Host

	pubSub    *msg.PubSub
	msgCenter *msg.MsgCenter
}

func NewNode(cfg util.Config) (Node, error) {
	node := Node{}
	node.ctx = context.Background()

	privKey, err := crypto.GenPrivKey()
	if err != nil {
		return Node{}, err
	}

	node.PrivKey = privKey
	node.PubKey = privKey.PubKey()
	node.Address = node.PubKey.Address()

	err = node.NewHost(cfg.ListenAddrs)
	if err != nil {
		return Node{}, err
	}

	err = node.NewPubSub()
	if err != nil {
		return Node{}, err
	}

	msgCenter, err := msg.NewMsgCenter(node.ctx, node.pubSub)
	if err != nil {
		return Node{}, err
	}

	node.msgCenter = msgCenter

	return node, nil
}

func (n *Node) Close() error {
	return n.host.Close()
}

func (n *Node) GetHostID() msg.ID {
	return n.host.ID()
}

func (n *Node) GetMsgCenter() *msg.MsgCenter {
	return n.msgCenter
}

func (n *Node) GetPeers() []msg.ID {
	return n.host.Network().Peers()
}

func (n *Node) Info() {
	if n.host == nil {
		return
	}

	fmt.Println("address:", n.Address)
	fmt.Println("host ID:", n.host.ID().Pretty())
	fmt.Println("host addrs:", n.host.Addrs())

	fmt.Printf("./petit-chat --bootstrap %s/p2p/%s\n",
		n.host.Addrs()[0],
		n.host.ID(),
	)
}
