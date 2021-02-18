package util

import (
	"flag"

	"github.com/h0n9/petit-chat/crypto"
)

type Config struct {
	BootstrapNodes crypto.Addrs
	ListenAddrs    crypto.Addrs
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) ParseFlags() error {
	flag.Var(&cfg.ListenAddrs, "listen", "addresses to listen from")
	flag.Var(&cfg.BootstrapNodes, "bootstrap", "bootstrap nodes")
	flag.Parse()

	return nil
}
