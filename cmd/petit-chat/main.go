package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/h0n9/petit-chat/cmd/petit-chat/cmd"
	"github.com/h0n9/petit-chat/msg"
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
	ctx := context.Background()
	node, err := p2p.NewNode(ctx, cfg)
	if err != nil {
		panic(err)
	}

	// handle signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		err = node.Close()
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}()

	err = node.DiscoverPeers(cfg.BootstrapNodes)
	if err != nil {
		panic(err)
	}

	hostPeer := msg.NewPeer(node.GetHostID(), "defaultHost")
	msgCenter, err := msg.NewCenter(ctx, node.GetPubSub(), hostPeer)
	if err != nil {
		panic(err)
	}

	err = node.SetCenter(hostPeer.GetNickname(), msgCenter)
	if err != nil {
		panic(err)
	}

	// CLI
	prompt := cmd.NewRootCmd(&node, hostPeer)
	err = prompt.Run()
	if err != nil {
		panic(err)
	}

	sigs <- syscall.SIGTERM
}

func init() {
	err := cfg.ParseFlags()
	if err != nil {
		panic(err)
	}
}
