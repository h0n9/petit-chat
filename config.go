package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
)

type Config struct {
	Context        context.Context
	BootstrapNodes Addrs
	ListenAddrs    Addrs
}

func (cfg *Config) parseFlags() error {
	flag.Var(&cfg.ListenAddrs, "listen", "addresses to listen from")
	flag.Var(&cfg.BootstrapNodes, "bootstrap", "bootstrap nodes")
	flag.Parse()

	if len(cfg.ListenAddrs) == 0 {
		randPort := rand.Intn(MaxListenPort-MinListenPort) + MinListenPort
		addr, err := NewMultiAddr(
			fmt.Sprintf("%s/%d/%s",
				DefaultListenAddr,
				randPort,
				TransportProtocol,
			),
		)
		if err != nil {
			return err
		}

		cfg.ListenAddrs = append(cfg.ListenAddrs, addr)
	}

	return nil
}
