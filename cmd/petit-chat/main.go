package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/h0n9/petit-chat/client"
	"github.com/h0n9/petit-chat/cmd/petit-chat/cmd"
	"github.com/h0n9/petit-chat/util"
)

// Expected implementation as follows
// - transport: QUIC
// - security: (expected to get supported by QUIC)
// - multiplex: (supported by QUIC)
// - peer discovery: Kademlia DHT
// - pubish subscribe: GossipSub

func main() {
	var cfg = util.Config{}
	err := cfg.ParseFlags()
	if err != nil {
		panic(err)
	}

	// init node
	ctx := context.Background()
	cli, err := client.NewClient(ctx, cfg)
	if err != nil {
		panic(err)
	}

	// handle signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		err = cli.Close()
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}()

	err = cli.DiscoverPeers()
	if err != nil {
		panic(err)
	}

	// CLI
	prompt := cmd.NewRootCmd(cli)
	err = prompt.Run()
	if err != nil {
		panic(err)
	}

	sigs <- syscall.SIGTERM
}
