package p2p

import (
	"fmt"

	libp2p "github.com/libp2p/go-libp2p"
	quic "github.com/libp2p/go-libp2p-quic-transport"

	"github.com/h0n9/petit-chat/crypto"
)

func (n *Node) NewHost(listenAddrs crypto.Addrs) error {
	if len(listenAddrs) == 0 {
		randPort := genRandPortNum()
		addr, err := crypto.NewMultiAddr(
			fmt.Sprintf("%s/%d/%s",
				DefaultListenAddr,
				randPort,
				TransportProtocol,
			),
		)
		if err != nil {
			return err
		}

		listenAddrs = append(listenAddrs, addr)
	}

	privKeyP2P, err := n.PrivKey.ToECDSAP2P()
	if err != nil {
		return err
	}

	tpt, err := quic.NewTransport(privKeyP2P)
	if err != nil {
		return err
	}

	host, err := libp2p.New(
		n.ctx,
		libp2p.ListenAddrs(listenAddrs...),
		libp2p.Identity(privKeyP2P),
		libp2p.Transport(tpt),
		libp2p.DefaultSecurity,
	)
	if err != nil {
		return err
	}

	n.host = host

	return nil
}
