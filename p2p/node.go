package p2p

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/libp2p/go-libp2p-core/crypto"

	"github.com/h0n9/petit-chat/net"
)

type Node struct {
	PrivKey crypto.PrivKey
	PubKey  crypto.PubKey

	Host net.Host
}

func NewNode(cfg Config) (Node, error) {
	node := Node{}
	privKey, pubKey, err := crypto.GenerateECDSAKeyPairWithCurve(
		elliptic.P256(),
		rand.Reader,
	)
	if err != nil {
		return Node{}, nil
	}

	node.PrivKey = privKey
	node.PubKey = pubKey

	host, err := net.NewHost(cfg.Context, node.PrivKey, cfg.ListenAddrs)
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

	fmt.Println("host ID:", n.Host.ID().Pretty())
	fmt.Println("host addrs:", n.Host.Addrs())

	fmt.Printf("%s/p2p/%s\n", n.Host.Addrs()[0], n.Host.ID())
}
