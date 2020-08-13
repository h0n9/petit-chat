package p2p

import (
	"context"
	"fmt"

	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type Node struct {
	ctx context.Context

	PrivKey crypto.PrivKey
	PubKey  crypto.PubKey
	Address crypto.Addr

	host Host

	pubSub *PubSub
}

func NewNode(ctx context.Context, cfg util.Config) (*Node, error) {
	privKey, err := crypto.GenPrivKey()
	if err != nil {
		return nil, err
	}

	node := Node{
		ctx:     ctx,
		PrivKey: privKey,
		PubKey:  privKey.PubKey(),
		Address: privKey.PubKey().Address(),
	}

	err = node.NewHost(cfg.ListenAddrs)
	if err != nil {
		return nil, err
	}

	err = node.NewPubSub()
	if err != nil {
		return nil, err
	}

	return &node, nil
}

func (n *Node) Close() error {
	return n.host.Close()
}

func (n *Node) GetHostID() types.ID {
	return n.host.ID()
}

func (n *Node) GetPeers() []types.ID {
	return n.host.Network().Peers()
}

func (n *Node) GetPubSub() *PubSub {
	return n.pubSub
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
