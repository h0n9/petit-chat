package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/h0n9/petit-chat/cmd/petit-chat/cmd"
	"github.com/h0n9/petit-chat/p2p"
	"github.com/h0n9/petit-chat/util"
)

// Expected implementation as follows
// - transport: QUIC
// - security: (expected to get supported by QUIC)
// - multiplex: (supported by QUIC)
// - peer discovery: Kademlia DHT
// - pubish subscribe: GossipSub

// global variables
var cfg util.Config

func main() {
	// init node
	node, err := p2p.NewNode(cfg)
	if err != nil {
		panic(err)
	}

	node.SetStreamHandler()

	err = node.DiscoverPeers(cfg.BootstrapNodes)
	if err != nil {
		panic(err)
	}

	// handle signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// CLI
	prompt := cmd.NewRootCmd(&node)
	err = prompt.Run()
	if err != nil {
		panic(err)
	}

	<-c

	err = node.Close()
	if err != nil {
		panic(err)
	}
}

func init() {
	err := cfg.ParseFlags()
	if err != nil {
		panic(err)
	}
}
