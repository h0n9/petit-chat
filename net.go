package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p-core/network"
)

// Multi Address used on QUIC protocol is formed as follows:
// ex) /ip4/0.0.0.0/udp/61881/quic

const (
	TransportProtocol = "quic"
	ProtocolID        = "/petit-chat/1.0.0"

	DefaultListenAddr = "/ip4/0.0.0.0/udp"
	MinListenPort     = 49152
	MaxListenPort     = 65535
)

func handleStream(stream network.Stream) {
	fmt.Println("new stream")

	// init buffer stream for non blocking read & write
	rw := bufio.NewReadWriter(
		bufio.NewReader(stream),
		bufio.NewWriter(stream),
	)

	// go routine for read & write data
	go readData(rw)
	go writeData(rw)
}

func readData(rw *bufio.ReadWriter) {
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
}

func writeData(rw *bufio.ReadWriter) {
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
}
