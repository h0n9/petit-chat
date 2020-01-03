package main

import (
	"fmt"

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

	node.SetStreamHandler()

	err = node.DiscoverPeers(cfg.BootstrapNodes)
	if err != nil {
		panic(err)
	}

	prompt := util.NewCmd("petit-chat", "entry point for petit-chat", nil,
		util.NewCmd("list", "list of subscribing topics",
			func(input string) error {
				for _, s := range node.GetSubs() {
					fmt.Printf("%s\n", s.Topic())
					for _, p := range node.GetPeers(s.Topic()) {
						fmt.Printf("  - %s\n", p)
					}
				}
				return nil
			},
		),
		util.NewCmd("pub", "publish to topic",
			func(input string) error {
				err := node.Publish("hello", []byte("wow"))
				if err != nil {
					return err
				}
				return nil
			},
		),
		util.NewCmd("sub", "subscribe to topic",
			func(input string) error {
				err := node.Subscribe("hello")
				if err != nil {
					return err
				}
				return nil
			},
		),
		util.NewCmd("unsub", "unsubscribe topic",
			func(input string) error {
				err := node.Unsubscribe("hello")
				if err != nil {
					return err
				}
				return nil
			},
		),
		util.NewCmd("test", "several cmds", nil,
			util.NewCmd("hello", "hello world",
				func(input string) error {
					fmt.Println("hello world")
					return nil
				},
			),
			util.NewCmd("other", "good", nil),
		),
	)

	prompt.Run()
}

func init() {
	err := cfg.ParseFlags()
	if err != nil {
		panic(err)
	}
}
