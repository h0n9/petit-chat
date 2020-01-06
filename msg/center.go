package msg

import (
	"crypto/sha256"

	"github.com/h0n9/petit-chat/code"
)

type MsgCenter struct {
	pubsub *PubSub

	Peers    map[Peer][]*MsgBox
	MsgBoxes map[string]*MsgBox
}

func NewMsgCenter(pubsub *PubSub) (*MsgCenter, error) {
	if pubsub == nil {
		return nil, code.ImproperPubSub
	}

	return &MsgCenter{
		pubsub:   pubsub,
		MsgBoxes: map[string]*MsgBox{},
	}, nil
}

func (mc *MsgCenter) SendMsg(data []byte, dests ...Peer) error {
	topic := genTopic(data, dests...)

	msgBox, exist := mc.MsgBoxes[topic]
	if !exist {
		msgBox = NewMsgBox(topic, dests)
	}

	err := mc.pubsub.Publish(topic, data)
	if err != nil {
		return err
	}

	return nil
}

func (mc *MsgCenter) GetMsgBoxes() map[string]*MsgBox {
	return mc.MsgBoxes
}

func (mc *MsgCenter) GetPeers() []Peer {
	peers := make([]Peer, 0, len(mc.Peers))

	for key := range mc.Peers {
		peers = append(peers, key)
	}

	return peers
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

func genTopic(init []byte, dests ...Peer) string {
	clay := make([]byte, 0, 50*len(dests)+len(init))

	for _, dest := range dests {
		mb, _ := dest.ID.MarshalBinary()
		clay = append(clay, mb...)
	}

	clay = append(clay, init...)
	hash := sha256.Sum256(clay)

	return string(hash[:])
}
