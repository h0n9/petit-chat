package net

import (
	"context"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func DiscoverPeers(ctx context.Context, h Host, bsNodes Addrs) error {
	// init peer discovery alg.
	peerDiscovery, err := dht.New(ctx, h)
	if err != nil {
		return err
	}

	// bootstrap peer discovery
	err = peerDiscovery.Bootstrap(ctx)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, bsn := range bsNodes {
		peerInfo, err := peer.AddrInfoFromP2pAddr(bsn)
		if err != nil {
			panic(err)
		}

		wg.Add(1)

		go func() {
			defer wg.Done()
			err = h.Connect(ctx, *peerInfo)
			if err != nil {
				panic(err)
			}

			fmt.Println("connected to:", *peerInfo)
		}()

	}
	wg.Wait()

	// advertise rendez-vous annoucement
	routingDiscovery := discovery.NewRoutingDiscovery(peerDiscovery)
	discovery.Advertise(ctx, routingDiscovery, RendezVous)

	peers, err := routingDiscovery.FindPeers(ctx, RendezVous)
	if err != nil {
		return err
	}

	for peer := range peers {
		if peer.ID == h.ID() {
			continue
		}

		err = h.Connect(ctx, peer)
		if err != nil {
			fmt.Println("failed to connect to:", peer)
			continue
		}

		fmt.Println("connected to:", peer)
	}

	return nil
}
