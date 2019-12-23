package main

import (
	"context"

	"github.com/h0n9/petit-chat/net"
	"github.com/h0n9/petit-chat/p2p"
)

// Expected implementation as follows
// - transport: QUIC
// - security: (expected to get supported by QUIC)
// - multiplex: (supported by QUIC)
// - peer discovery: Kademlia DHT

// global variables
var (
	cfg  p2p.Config
	node p2p.Node
)

func main() {
	node, err := p2p.NewNode(cfg)
	if err != nil {
		panic(err)
	}

	node.Info()

	net.SetStreamHandler(node.Host)

	err = net.DiscoverPeers(cfg.Context, node.Host, cfg.BootstrapNodes)
	if err != nil {
		panic(err)
	}

	// TODO: pub/sub nodes

	// to keep the app alive
	select {}
}

func init() {
	err := cfg.ParseFlags()
	if err != nil {
		panic(err)
	}

	cfg.Context = context.Background()
}
