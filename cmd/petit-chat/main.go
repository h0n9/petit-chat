package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

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
	node, err := p2p.NewNode(cfg)
	if err != nil {
		panic(err)
	}

	node.Info()

	node.SetStreamHandler()

	err = node.DiscoverPeers(cfg.BootstrapNodes)
	if err != nil {
		panic(err)
	}

	prompt := util.NewCmd("petit-chat", "entry point for petit-chat", nil,
		util.NewCmd("list", "list of subscribing topics",
			func(reader *bufio.Reader) error {
				fmt.Println("list of message boxes")
				listOfMsgBoxes(node)
				fmt.Println("\nlist of peers (network)")
				listOfPeers(node)
				return nil
			},
		),
		util.NewCmd("send", "send a message to peers",
			func(reader *bufio.Reader) error {
				msgCenter := node.GetMsgCenter()
				peers := listOfPeers(node)

				from := msg.Peer{ID: node.GetHostID()}

				fmt.Println("choose peers to send a message")
				data, err := util.GetInput(reader)
				if err != nil {
					return err
				}

				nums := strings.Split(data, ",")
				if len(nums) == 0 {
					return fmt.Errorf("type peers")
				}

				tos := make([]msg.Peer, 0, len(nums))
				for _, num := range nums {
					i, err := strconv.Atoi(num)
					if err != nil {
						return err
					}

					if i < 1 || i > len(nums) {
						return fmt.Errorf("unavailable input as peer index: %d", i)
					}

					to := msg.Peer{ID: peers[i-1]}
					tos = append(tos, to)
				}

				fmt.Println("type a message to send")
				data, err = util.GetInput(reader)
				if err != nil {
					return err
				}

				err = msgCenter.SendMsg([]byte(data), from, tos)
				if err != nil {
					return err
				}

				return nil
			},
		),
		util.NewCmd("test", "several cmds", nil,
			util.NewCmd("hello", "hello world",
				func(reader *bufio.Reader) error {
					fmt.Println("hello world")
					return nil
				},
			),
			util.NewCmd("other", "good", nil),
		),
	)

	err = prompt.Run()
	if err != nil {
		panic(err)
	}
}

func listOfMsgBoxes(node p2p.Node) {
	msgCenter := node.GetMsgCenter()

	for topic, msgBox := range msgCenter.GetMsgBoxes() {
		fmt.Printf("%s\n", topic)
		for _, p := range msgBox.GetPeers() {
			fmt.Printf("  - %s\n", p)
		}
	}
}

func listOfPeers(node p2p.Node) []msg.ID {
	peers := node.GetPeers()

	for i, peer := range peers {
		fmt.Printf("%d. %s\n", i+1, peer)
	}

	return peers
}

func init() {
	err := cfg.ParseFlags()
	if err != nil {
		panic(err)
	}
}
