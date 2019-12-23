package net

import (
	"bufio"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p-core/network"
)

func SetStreamHandler(h Host) {
	h.SetStreamHandler(ProtocolID, handleStream)
}

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
