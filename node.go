package main

import (
	"crypto/elliptic"
	"crypto/rand"

	"github.com/libp2p/go-libp2p-core/crypto"
)

type Node struct {
	PrivKey crypto.PrivKey
	PubKey  crypto.PubKey
}

func NewNode() (Node, error) {
	node := Node{}
	privKey, pubKey, err := crypto.GenerateECDSAKeyPairWithCurve(elliptic.P256(), rand.Reader)
	if err != nil {
		return Node{}, nil
	}

	node.PrivKey = privKey
	node.PubKey = pubKey

	return node, nil
}
