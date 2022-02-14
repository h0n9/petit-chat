package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/h0n9/petit-chat/client"
	"github.com/h0n9/petit-chat/cmd/petit-chat/cmd"
	"github.com/h0n9/petit-chat/server"
	"github.com/h0n9/petit-chat/util"
)

// Expected implementation as follows
// - transport: QUIC
// - security: (expected to get supported by QUIC)
// - multiplex: (supported by QUIC)
// - peer discovery: Kademlia DHT
// - pubish subscribe: GossipSub

func main() {
	cfg := util.NewConfig()
	err := cfg.ParseFlags()
	if err != nil {
		panic(err)
	}

	// init server
	ctx := context.Background()
	svr, err := server.NewServer(ctx, cfg)
	if err != nil {
		panic(err)
	}

	// init client
	cli, err := client.NewClient(svr)
	if err != nil {
		panic(err)
	}

	// handle signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		err = svr.Close()
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}()

	err = svr.DiscoverPeers()
	if err != nil {
		panic(err)
	}

	// CLI
	prompt := cmd.NewRootCmd(svr, cli)
	err = prompt.Run()
	if err != nil {
		panic(err)
	}

	sigs <- syscall.SIGTERM
}
