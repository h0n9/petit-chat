package main

import (
	"github.com/h0n9/petit-chat/p2p"
	"github.com/h0n9/petit-chat/util"
)

// Expected implementation as follows
// - transport: QUIC
// - security: (expected to get supported by QUIC)
// - multiplex: (supported by QUIC)
// - peer discovery: Kademlia DHT

// global variables
var (
	cfg  util.Config
	node p2p.Node
)

func main() {
	node, err := p2p.NewNode(cfg)
	if err != nil {
		panic(err)
	}

	node.Info()

	p2p.SetStreamHandler(node.Host)

	err = node.DiscoverPeers(cfg.BootstrapNodes)
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
}
