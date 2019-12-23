package main

import (
	"context"
	"fmt"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	quic "github.com/libp2p/go-libp2p-quic-transport"

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
	n, err := p2p.NewNode()
	if err != nil {
		panic(err)
	}

	node = n

	quicTpt, err := quic.NewTransport(node.PrivKey)
	if err != nil {
		panic(err)
	}

	host, err := libp2p.New(
		cfg.Context,
		libp2p.ListenAddrs(cfg.ListenAddrs...),
		libp2p.Identity(node.PrivKey),
		libp2p.Transport(quicTpt),
		libp2p.DefaultSecurity,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("host ID:", host.ID().Pretty())
	fmt.Println("host addrs:", host.Addrs())

	fmt.Printf("%s/p2p/%s\n", host.Addrs()[0], host.ID().Pretty())

	host.SetStreamHandler(protocol.ID(net.ProtocolID), net.HandleStream)

	// init peer discovery alg.
	peerDiscovery, err := dht.New(cfg.Context, host)
	if err != nil {
		panic(err)
	}

	// bootstrap peer discovery
	err = peerDiscovery.Bootstrap(cfg.Context)
	if err != nil {
		panic(err)
	}

	// TODO: connect to bootstrap nodes
	var wg sync.WaitGroup
	for _, bsn := range cfg.BootstrapNodes {
		peerInfo, err := peer.AddrInfoFromP2pAddr(bsn)
		if err != nil {
			panic(err)
		}

		wg.Add(1)

		go func() {
			defer wg.Done()
			err = host.Connect(cfg.Context, *peerInfo)
			if err != nil {
				panic(err)
			}

			fmt.Println("connected to:", *peerInfo)
		}()

	}
	wg.Wait()

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
