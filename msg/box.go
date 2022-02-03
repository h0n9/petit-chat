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
	ctx       context.Context
	chCapsule chan *Capsule

	hostID types.ID
	topic  *types.Topic
	sub    *types.Sub

	state *types.State
	store *CapsuleStore
}

func NewBox(ctx context.Context, topic *types.Topic, public bool, hostID types.ID) (*Box, error) {
	return &Box{
		ctx:       ctx,
		chCapsule: make(chan *Capsule, 1),

		hostID: hostID,
		topic:  topic,
		sub:    nil,

		state: types.NewState(public),
		store: NewCapsuleStore(),
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

func (box *Box) Publish(capsule *Capsule) error {
	data, err := capsule.Bytes()
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
		capsule, err := NewCapsuleFromBytes(received.GetData())
		if err != nil {
			fmt.Println(err)
			continue
		}
		// TODO: add constraints to capsule
		err = capsule.Check()
		if err != nil {
			fmt.Println(err)
			continue
		}

		box.chCapsule <- capsule

		_, err = box.append(capsule)
		if err != nil {
			fmt.Println(err)
			continue
		}

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

func (box *Box) GetChCapsule() chan *Capsule {
	return box.chCapsule
}

func (box *Box) GetCapsules() []*Capsule {
	return box.store.GetCapsules()
}

func (box *Box) GetCapsule(hash types.Hash) *Capsule {
	return box.store.GetCapsule(hash)
}

// func (box *Box) GetUnreadMsgs() []*Capsule {
// 	Capsules := []*Capsule{}
// 	readUntilIndex := box.state.GetReadUntilIndex()
// 	if readUntilIndex+1 < uint64(len(box.store.Capsules)) {
// 		Capsules = append(Capsules, box.store.Capsules[readUntilIndex+1:]...)
// 	}
// 	box.state.SetReadUntilIndex(uint64(len(box.store.Capsules) - 1))
// 	return Capsules
// }

func (box *Box) GetAuth() *types.Auth {
	return box.state.GetAuth()
}

func (box *Box) append(capsule *Capsule) (types.Index, error) {
	return box.store.Append(capsule)
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
