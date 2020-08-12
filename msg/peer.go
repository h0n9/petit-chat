package msg

import (
	"github.com/libp2p/go-libp2p-core/peer"
)

type ID = peer.ID

type Peer struct {
	id       ID
	nickname string
}

func NewPeer(id ID, nickname string) Peer {
	return Peer{
		id:       id,
		nickname: nickname,
	}
}

func (p *Peer) GetID() ID {
	return p.id
}

func (p *Peer) GetNickname() string {
	return p.nickname
}
