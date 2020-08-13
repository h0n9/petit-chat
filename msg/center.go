package msg

import (
	"github.com/h0n9/petit-chat/code"
)

type Center struct {
	// TODO: better way to manage topics with peerList
	msgBoxes map[string]*Box
}

func NewCenter() (*Center, error) {
	return &Center{msgBoxes: make(map[string]*Box)}, nil
}

func (mc *Center) CreateBox(topic string) (*Box, error) {
	_, exist := mc.getBox(topic)
	if exist {
		return nil, code.AlreadyExistingTopic
	}

	msgBox, err := NewBox()
	if err != nil {
		return nil, err
	}

	err = mc.add(topic, msgBox)
	if err != nil {
		return nil, err
	}

	return msgBox, nil
}

func (mc *Center) LeaveBox(topic string) error {
	_, exist := mc.getBox(topic)
	if !exist {
		return code.NonExistingTopic
	}

	delete(mc.msgBoxes, topic)

	return nil
}

func (mc *Center) GetBoxes() map[string]*Box {
	return mc.msgBoxes
}

func (mc *Center) GetBox(topic string) (*Box, error) {
	msgBox, exist := mc.getBox(topic)
	if !exist {
		return nil, code.NonExistingTopic
	}

	return msgBox, nil
}

func (mc *Center) getBox(topic string) (*Box, bool) {
	msgBox, exist := mc.msgBoxes[topic]
	return msgBox, exist
}

func (mc *Center) add(topic string, msgBox *Box) error {
	_, exist := mc.getBox(topic)
	if exist {
		return code.AlreadyExistingTopic
	}

	mc.msgBoxes[topic] = msgBox

	return nil
}

func (mc *Center) remove(topic string) error {
	_, exist := mc.getBox(topic)
	if !exist {
		return code.NonExistingTopic
	}

	delete(mc.msgBoxes, topic)

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
