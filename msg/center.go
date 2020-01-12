package msg

import (
	"context"
	"crypto/sha256"
	"sort"

	"github.com/h0n9/petit-chat/code"
)

type MsgCenter struct {
	ctx    context.Context
	pubsub *PubSub

	// TODO: better way to manage topics with peerList
	Peers    []Peer
	MsgBoxes map[string]*MsgBox
}

func NewMsgCenter(ctx context.Context, pubsub *PubSub, peers ...Peer,
) (*MsgCenter, error) {
	if pubsub == nil {
		return nil, code.ImproperPubSub
	}

	return &MsgCenter{
		ctx:    ctx,
		pubsub: pubsub,

		Peers:    peers,
		MsgBoxes: make(map[string]*MsgBox),
	}, nil
}

func (mc *MsgCenter) SendMsg(data []byte, from Peer, tos []Peer) error {
	topic := genTopic(data, append(tos, from))

	msgBox, exist := mc.MsgBoxes[topic]
	if !exist {
		sub, err := mc.pubsub.Subscribe(topic)
		if err != nil {
			return err
		}

		msgBox, err := NewMsgBox(mc.ctx, sub, from, "", tos...)
		if err != nil {
			return err
		}

		mc.add(topic, msgBox)
	}

	msg := NewMsg(data, from, tos)
	msgJSON, err := msg.MarshalJSON()
	if err != nil {
		return err
	}

	err = mc.pubsub.Publish(topic, msgJSON)
	if err != nil {
		return err
	}

	err = msgBox.Append(msg)
	if err != nil {
		return err
	}

	return nil
}

func (mc *MsgCenter) GetMsgBoxes() map[string]*MsgBox {
	return mc.MsgBoxes
}

func (mc *MsgCenter) GetPeers() []Peer {
	return mc.Peers
}

func (mc *MsgCenter) add(topic string, msgBox *MsgBox) error {
	_, exist := mc.MsgBoxes[topic]
	if exist {
		return code.AlreadyExistingTopic
	}

	mc.MsgBoxes[topic] = msgBox

	return nil
}

func (mc *MsgCenter) remove(topic string) error {
	_, exist := mc.MsgBoxes[topic]
	if !exist {
		return code.NonExistingTopic
	}

	delete(mc.MsgBoxes, topic)

	return nil
}

func genTopic(init []byte, peers []Peer) string {
	// TODO: compact way to allocate memory for clay variable
	clay := make([]byte, 0, 50*(len(peers))+len(init))

	sort.Slice(peers, func(i, j int) bool {
		return peers[i].ID > peers[j].ID
	})

	for _, peer := range peers {
		b, _ := peer.ID.MarshalBinary()
		clay = append(clay, b...)
	}

	clay = append(clay, init...)
	hash := sha256.Sum256(clay)

	return string(hash[:])
}
