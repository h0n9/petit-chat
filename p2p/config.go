package p2p

import (
	"context"
	"flag"
	"fmt"
	"math/rand"

	"github.com/h0n9/petit-chat/net"
)

type Config struct {
	Context        context.Context
	BootstrapNodes net.Addrs
	ListenAddrs    net.Addrs
}

func (cfg *Config) ParseFlags() error {
	flag.Var(&cfg.ListenAddrs, "listen", "addresses to listen from")
	flag.Var(&cfg.BootstrapNodes, "bootstrap", "bootstrap nodes")
	flag.Parse()

	if len(cfg.ListenAddrs) == 0 {
		randPort := rand.Intn(net.MaxListenPort-net.MinListenPort) + net.MinListenPort
		addr, err := net.NewMultiAddr(
			fmt.Sprintf("%s/%d/%s",
				net.DefaultListenAddr,
				randPort,
				net.TransportProtocol,
			),
		)
		if err != nil {
			return err
		}

		cfg.ListenAddrs = append(cfg.ListenAddrs, addr)
	}

	return nil
}
