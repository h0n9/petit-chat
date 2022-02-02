package msg

import (
	"context"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
)

type Center struct {
	// TODO: better way to manage topics with peerList
	ctx      context.Context
	myID     types.ID
	msgBoxes map[string]*Box
}

func NewCenter(ctx context.Context, myID types.ID) (*Center, error) {
	return &Center{
		ctx:      ctx,
		myID:     myID,
		msgBoxes: make(map[string]*Box),
	}, nil
}

func (mc *Center) CreateBox(tp *types.Topic, pub bool, pk *crypto.PrivKey, p *types.Persona) (*Box, error) {
	tStr := tp.String()
	_, exist := mc.getBox(tStr)
	if exist {
		return nil, code.AlreadyExistingTopic
	}

	msgBox, err := NewBox(mc.ctx, tp, pub, mc.myID)
	if err != nil {
		return nil, err
	}

	err = mc.add(tStr, msgBox)
	if err != nil {
		return nil, err
	}

	return msgBox, nil
}

func (mc *Center) LeaveBox(topicStr string) error {
	box, exist := mc.getBox(topicStr)
	if !exist {
		return code.NonExistingTopic
	}

	if box.Subscribing() {
		err := box.Close()
		if err != nil {
			return err
		}
	}

	err := mc.remove(topicStr)
	if err != nil {
		return err
	}

	return nil
}

func (mc *Center) GetBoxes() map[string]*Box {
	return mc.msgBoxes
}

func (mc *Center) GetBox(topicStr string) (*Box, bool) {
	return mc.getBox(topicStr)
}

func (mc *Center) getBox(topicStr string) (*Box, bool) {
	msgBox, exist := mc.msgBoxes[topicStr]
	return msgBox, exist
}

func (mc *Center) add(topicStr string, msgBox *Box) error {
	_, exist := mc.getBox(topicStr)
	if exist {
		return code.AlreadyExistingTopic
	}

	mc.msgBoxes[topicStr] = msgBox

	return nil
}

func (mc *Center) remove(topicStr string) error {
	_, exist := mc.getBox(topicStr)
	if !exist {
		return code.NonExistingTopic
	}

	delete(mc.msgBoxes, topicStr)

	return nil
}

// func genTopic(init []byte, peers []Peer) string {
// 	// TODO: compact way to allocate memory for clay variable
// 	clay := make([]byte, 0, 50*(len(peers))+len(init))
//
// 	sort.Slice(peers, func(i, j int) bool {
// 		return peers[i].GetID() > peers[j].GetID()
// 	})
//
// 	for _, peer := range peers {
// 		b, _ := peer.GetID().MarshalBinary()
// 		clay = append(clay, b...)
// 	}
//
// 	clay = append(clay, init...)
// 	hash := sha256.Sum256(clay)
//
// 	return string(hash[:])
// }
