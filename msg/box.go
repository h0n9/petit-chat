package msg

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

// Box refers to a chat room
type Box struct {
	ctx          context.Context
	chMsgCapsule chan *MsgCapsule

	topic *types.Topic
	sub   *types.Sub

	hostID types.ID
	state  *types.State

	store struct {
		msgs      []*Msg              // TODO: limit the size of msgs slice
		msgHashes map[types.Hash]*Msg // TODO: limit the size of msgHashes map
	}
}

func NewBox(ctx context.Context, topic *types.Topic, public bool, hostID types.ID) (*Box, error) {
	return &Box{
		ctx:          ctx,
		chMsgCapsule: make(chan *MsgCapsule, 1),

		topic: topic,
		sub:   nil,

		hostID: hostID,
		state:  types.NewState(public),

		store: struct {
			msgs      []*Msg
			msgHashes map[types.Hash]*Msg
		}{
			msgs:      make([]*Msg, 0),
			msgHashes: make(map[types.Hash]*Msg),
		},
	}, nil
	// err = box.join(hostPersona)
	// if err != nil {
	// 	return nil, err
	// }
	// err = grant(box.state.GetAuth(), hostPersona.Address, true, true, true)
	// if err != nil {
	// 	return nil, err
	// }
	// msg := NewMsgHelloSyn(&box, types.EmptyHash, hostPersona)
	// err = box.Publish(msg, false)
	// if err != nil {
	// 	return nil, err
	// }
	// return &box, nil
}

func (box *Box) Publish(msgCapsule *MsgCapsule) error {
	data, err := msgCapsule.Bytes()
	if err != nil {
		return err
	}
	err = box.topic.Publish(box.ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (box *Box) Subscribe() error {
	if box.sub != nil {
		return code.AlreadySubscribingTopic
	}

	sub, err := box.topic.Subscribe()
	if err != nil {
		return err
	}
	box.sub = sub

	for {
		received, err := sub.Next(box.ctx)
		if err != nil {
			// TODO: replace fmt.Println() to logger.Println()
			fmt.Println(err)
			continue
		}
		msgCapsule, err := NewMsgCapsuleFromBytes(received.GetData())
		if err != nil {
			return err
		}
		// TODO: add constraints to msgCapsule
		box.chMsgCapsule <- msgCapsule
		// msg, err := box.Decapsulate(data)
		// if err != nil {
		// 	fmt.Println(err)
		// 	continue
		// }
		// err = box.Verify(msg)
		// if err != nil {
		// 	fmt.Println(err)
		// 	continue
		// }
		// eos, err := handler(box, msg)
		// if err != nil {
		// 	// TODO: replace fmt.Println() to logger.Println()
		// 	fmt.Println(err)
		// 	continue
		// }

		// // eos shoud be the only way to break for loop
		// if eos {
		// 	box.sub.Cancel()
		// 	err = box.topic.Close()
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// 	break
		// }
	}
}

func Hash(base Base) types.Hash {
	data, err := json.Marshal(base)
	if err != nil {
		return types.EmptyHash
	}
	return util.ToSHA256(data)
}

func (box *Box) GetPersonae() map[crypto.Addr]*types.Persona {
	return box.state.GetPersonae()
}

func (box *Box) getPersona(addr crypto.Addr) *types.Persona {
	return box.state.GetPersona(addr)
}

func (box *Box) GetPersona(addr crypto.Addr) *types.Persona {
	return box.getPersona(addr)
}

func (box *Box) GetHostID() types.ID {
	return box.hostID
}

func grant(auth *types.Auth, addr crypto.Addr, r, w, x bool) error {
	perm := types.NewPerm(r, w, x)
	err := auth.SetPerm(addr, perm)
	if err != nil {
		return err
	}
	return nil
}

func (box *Box) Grant(addr crypto.Addr, r, w, x bool) error {
	auth := box.state.GetAuth()
	newAuth, err := auth.Copy()
	if err != nil {
		return err
	}
	err = grant(newAuth, addr, r, w, x)
	if err != nil {
		return err
	}

	personae := box.state.GetPersonae()
	err = box.propagate(newAuth, personae)
	if err != nil {
		return err
	}

	return nil
}

func revoke(auth *types.Auth, addr crypto.Addr) error {
	err := auth.DeletePerm(addr)
	if err != nil {
		return err
	}
	return nil
}

func (box *Box) Revoke(addr crypto.Addr) error {
	auth := box.state.GetAuth()
	newAuth, err := auth.Copy()
	if err != nil {
		return err
	}
	err = revoke(newAuth, addr)
	if err != nil {
		return err
	}

	personae := box.state.GetPersonae()
	err = box.propagate(newAuth, personae)
	if err != nil {
		return err
	}

	return nil
}

func (box *Box) Close() error {
	// Announe EOS to others (application layer)
	// msg := NewMsgBye(box, types.EmptyHash, box.vault.GetPersona())
	// return box.Publish(msg, true)
	return nil
}

func (box *Box) Subscribing() bool {
	return box.sub != nil
}

func (box *Box) GetChMsgCapsule() chan *MsgCapsule {
	return box.chMsgCapsule
}

func (box *Box) GetMsgs() []*Msg {
	return box.store.msgs
}

func (box *Box) GetMsg(mh types.Hash) *Msg {
	return box.store.msgHashes[mh]
}

func (box *Box) GetUnreadMsgs() []*Msg {
	msgs := []*Msg{}
	readUntilIndex := box.state.GetReadUntilIndex()
	if readUntilIndex+1 < uint64(len(box.store.msgs)) {
		msgs = append(msgs, box.store.msgs[readUntilIndex+1:]...)
	}
	box.state.SetReadUntilIndex(uint64(len(box.store.msgs) - 1))
	return msgs
}

func (box *Box) GetAuth() *types.Auth {
	return box.state.GetAuth()
}

func (box *Box) append(msg *Msg) (types.Index, error) {
	hash := msg.GetHash()
	_, exist := box.store.msgHashes[hash]
	if exist {
		return 0, code.AlreadyAppendedMsg
	}

	timestamp := msg.GetTimestamp()
	latestTimestamp := box.state.GetLatestTimestamp()
	if latestTimestamp.Before(timestamp) {
		box.state.SetLatestTimestamp(timestamp)
	}

	box.store.msgs = append(box.store.msgs, msg)
	box.store.msgHashes[hash] = msg

	return types.Index(len(box.store.msgs) - 1), nil
}

func (box *Box) join(targetPersona *types.Persona) error {
	oldPersona := box.getPersona(targetPersona.Address)
	if oldPersona != nil {
		return nil // ignore even if existing
		// return code.ExistingPersonaInBox
	}
	err := targetPersona.Check()
	if err != nil {
		return err
	}
	box.state.SetPersona(targetPersona.Address, targetPersona)
	return nil
}

func (box *Box) leave(targetPersona *types.Persona) error {
	oldPersona := box.getPersona(targetPersona.Address)
	if oldPersona == nil {
		return code.NonExistingPersonaInBox
	}
	box.state.DeletePersona(targetPersona.Address)
	return nil
}

func (box *Box) propagate(auth *types.Auth, personae types.Personae) error {
	// msg := NewMsgUpdate(box, types.EmptyHash, auth, personae)
	// err := box.Publish(msg, true)
	// if err != nil {
	// 	return err
	// }
	return nil
}
