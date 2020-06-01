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
	host     Peer
	peers    []Peer
	msgBoxes map[string]*MsgBox
}

func NewMsgCenter(ctx context.Context, pubsub *PubSub, host Peer, peers ...Peer,
) (*MsgCenter, error) {
	if pubsub == nil {
		return nil, code.ImproperPubSub
	}

	return &MsgCenter{
		ctx:    ctx,
		pubsub: pubsub,

		host:     host,
		peers:    peers,
		msgBoxes: make(map[string]*MsgBox),
	}, nil
}

func (mc *MsgCenter) CreateMsgBox(topic string) error {
	err := mc.check(topic)
	if err != nil {
		return err
	}

	sub, err := mc.pubsub.Subscribe(topic)
	if err != nil {
		return err
	}

	msgBox, err := NewMsgBox(mc.ctx, sub, mc.host)
	if err != nil {
		return err
	}

	err = mc.add(topic, msgBox)
	if err != nil {
		return err
	}

	return nil
}

func (mc *MsgCenter) GetMsgBoxes() map[string]*MsgBox {
	return mc.msgBoxes
}

func (mc *MsgCenter) GetPeers() []Peer {
	return mc.peers
}

func (mc *MsgCenter) check(topic string) error {
	_, exist := mc.msgBoxes[topic]
	if exist {
		return code.AlreadyExistingTopic
	}

	return nil
}

func (mc *MsgCenter) add(topic string, msgBox *MsgBox) error {
	err := mc.check(topic)
	if err != nil {
		return err
	}

	mc.msgBoxes[topic] = msgBox

	return nil
}

func (mc *MsgCenter) remove(topic string) error {
	_, exist := mc.msgBoxes[topic]
	if !exist {
		return code.NonExistingTopic
	}

	delete(mc.msgBoxes, topic)

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
