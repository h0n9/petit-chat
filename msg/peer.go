package msg

import (
	"github.com/libp2p/go-libp2p-core/peer"
)

type ID = peer.ID

type Peer struct {
	ID       ID
	Nickname string
}
