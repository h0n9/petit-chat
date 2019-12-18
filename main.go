package main

import (
	"bufio"
	"context"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"os"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	quic "github.com/libp2p/go-libp2p-quic-transport"
)

// Expected implementation as follows
// - transport: QUIC
// - security: (expected to get supported by QUIC)
// - multiplex: (supported by QUIC)
// - peer discovery: Kademlia DHT

func handleStream(stream network.Stream) {
	fmt.Println("new stream")

	// init buffer stream for non blocking read & write
	rw := bufio.NewReadWriter(
		bufio.NewReader(stream),
		bufio.NewWriter(stream),
	)

	// go routine for read & write data
	go func(rw *bufio.ReadWriter) {
		for {
			data, err := rw.ReadString('\n')
			if err != nil {
				panic(err)
			}

			// ignore empty line
			if data == "" || data == "\n" {
				return
			}

			fmt.Println(data)
		}
	}(rw)
	go func(rw *bufio.ReadWriter) {
		// reader(stdin) -> data -> bufio -> stream

		// init reader reading data from stdin
		reader := bufio.NewReader(os.Stdin)

		for {
			fmt.Printf("> ")

			// reader -> data
			data, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}

			// data -> bufio
			_, err = rw.WriteString(data + "\n")
			if err != nil {
				panic(err)
			}

			// bufio -> stream (flush)
			err = rw.Flush()
			if err != nil {
				panic(err)
			}
		}
	}(rw)
}

func main() {
	// init of empty context
	ctx := context.Background()

	privKey, pubKey, err := crypto.GenerateECDSAKeyPairWithCurve(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	fmt.Println("private key:", privKey)
	fmt.Println("public key: ", pubKey)

	quicTpt, err := quic.NewTransport(privKey)
	if err != nil {
		panic(err)
	}

	host, err := libp2p.New(ctx, libp2p.Transport(quicTpt), libp2p.DefaultSecurity)
	if err != nil {
		panic(err)
	}

	fmt.Println("host ID:", host.ID())

	host.SetStreamHandler(protocol.ID("petit-chat"), handleStream)

	// init peer discovery alg.
	peerDiscovery, err := dht.New(ctx, host)
	if err != nil {
		panic(err)
	}

	// bootstrap peer discovery
	err = peerDiscovery.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}

	// TODO: connect to bootstrap nodes

	// to keep the app alive
	select {}
}
