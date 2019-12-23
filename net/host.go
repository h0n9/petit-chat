package net

import (
	"context"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	quic "github.com/libp2p/go-libp2p-quic-transport"
)

type Host = host.Host

func NewHost(ctx context.Context, privKey crypto.PrivKey, listenAddrs Addrs) (Host, error) {
	tpt, err := quic.NewTransport(privKey)
	if err != nil {
		panic(err)
	}

	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(listenAddrs...),
		libp2p.Identity(privKey),
		libp2p.Transport(tpt),
		libp2p.DefaultSecurity,
	)
	if err != nil {
		return nil, err
	}

	return host, nil
}
