package msg

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
)

// Box refers to a chat room
type Box struct {
	ctx       context.Context
	topic     *types.Topic
	sub       *types.Sub
	auth      *types.Auth
	secretKey *crypto.SecretKey

	myID            types.ID
	myPrivKey       *crypto.PrivKey
	myPersona       *types.Persona
	msgSubCh        chan *Msg
	latestTimestamp time.Time
	readUntilIndex  int

	personae  types.Personae
	msgs      []*Msg              // TODO: limit the size of msgs slice
	msgHashes map[types.Hash]*Msg // TODO: limit the size of msgHashes map
}

func NewBox(ctx context.Context, topic *types.Topic, public bool,
	myID types.ID, myPrivKey *crypto.PrivKey, myPersona *types.Persona) (*Box, error) {
	err := myPersona.Check()
	if err != nil {
		return nil, err
	}
	secretKey, err := crypto.GenSecretKey()
	if err != nil {
		return nil, err
	}
	box := Box{
		ctx:       ctx,
		topic:     topic,
		sub:       nil,
		auth:      types.NewAuth(public, make(map[crypto.Addr]types.Perm)),
		secretKey: secretKey,

		myID:            myID,
		myPrivKey:       myPrivKey,
		myPersona:       myPersona,
		msgSubCh:        nil,
		latestTimestamp: time.Now(),
		readUntilIndex:  0,

		personae:  make(types.Personae),
		msgs:      make([]*Msg, 0),
		msgHashes: make(map[types.Hash]*Msg),
	}
	err = box.join(myPersona)
	if err != nil {
		return nil, err
	}
	err = grant(box.auth, box.myPersona.Address, true, true, true)
	if err != nil {
		return nil, err
	}
	msg := NewMsg(box.myID, types.EmptyHash, &BodyHelloSyn{
		Persona: myPersona,
	})
	err = box.Publish(msg, TypeHelloSyn, false)
	if err != nil {
		return nil, err
	}
	return &box, nil
}

func (box *Box) Encapsulate(msg *Msg, msgType Type, encrypt bool) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	if encrypt {
		data, err = box.secretKey.Encrypt(data)
		if err != nil {
			return nil, err
		}
	}

	data, err = json.Marshal(&MsgCapsule{
		Encrypted: encrypt,
		Type:      msgType,
		Data:      data,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (box *Box) Decapsulate(data []byte) (*Msg, error) {
	msgCapsule := MsgCapsule{}
	msg := Msg{}

	err := json.Unmarshal(data, &msgCapsule)
	if err != nil {
		return nil, err
	}

	if msgCapsule.Encrypted {
		msgCapsule.Data, err = box.secretKey.Decrypt(msgCapsule.Data)
		if err != nil {
			return nil, err
		}
	}

	body := msgCapsule.Type.Body()
	if body == nil {
		return nil, code.UnknownMsgType
	}
	msg.Body = body

	err = json.Unmarshal(msgCapsule.Data, &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (box *Box) Publish(msg *Msg, msgType Type, encrypt bool) error {
	err := msg.Sign(box.myPrivKey)
	if err != nil {
		return err
	}
	data, err := box.Encapsulate(msg, msgType, encrypt)
	if err != nil {
		return err
	}
	err = box.topic.Publish(box.ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (box *Box) Subscribe(handler MsgHandler) error {
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
		data := received.GetData()
		msg, err := box.Decapsulate(data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = msg.Verify()
		if err != nil {
			fmt.Println(err)
			continue
		}
		eos, err := handler(box, msg)
		if err != nil {
			// TODO: replace fmt.Println() to logger.Println()
			fmt.Println(err)
			continue
		}

		// eos shoud be the only way to break for loop
		if eos {
			box.sub.Cancel()
			err = box.topic.Close()
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}

	return nil
}

func (box *Box) GetPersonae() map[crypto.Addr]*types.Persona {
	return box.personae
}

func (box *Box) getPersona(cAddr crypto.Addr) *types.Persona {
	persona, exist := box.personae[cAddr]
	if !exist {
		return nil
	}
	return persona
}

func (box *Box) GetPersona(cAddr crypto.Addr) *types.Persona {
	return box.getPersona(cAddr)
}

func (box *Box) GetMyID() types.ID {
	return box.myID
}

func (box *Box) GetMyPersona() *types.Persona {
	return box.myPersona
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
	newAuth, err := box.auth.Copy()
	if err != nil {
		return err
	}
	err = grant(newAuth, addr, r, w, x)
	if err != nil {
		return err
	}

	err = box.propagate(newAuth, box.personae)
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
	newAuth, err := box.auth.Copy()
	if err != nil {
		return err
	}
	err = revoke(newAuth, addr)
	if err != nil {
		return err
	}

	err = box.propagate(newAuth, box.personae)
	if err != nil {
		return err
	}

	return nil
}

func (box *Box) Close() error {
	// Announe EOS to others (application layer)
	msg := NewMsg(box.myID, types.EmptyHash, &BodyBye{
		Persona: box.myPersona,
	})
	return box.Publish(msg, TypeBye, true)
}

func (box *Box) Subscribing() bool {
	return box.sub != nil
}

func (box *Box) SetMsgSubCh(msgSubCh chan *Msg) {
	box.msgSubCh = msgSubCh
}

func (box *Box) GetSecretKey() *crypto.SecretKey {
	return box.secretKey
}

func (box *Box) GetMsgs() []*Msg {
	return box.msgs
}

func (box *Box) GetMsg(mh types.Hash) *Msg {
	return box.msgHashes[mh]
}

func (box *Box) GetUnreadMsgs() []*Msg {
	msgs := []*Msg{}
	if box.readUntilIndex+1 < len(box.msgs) {
		msgs = append(msgs, box.msgs[box.readUntilIndex+1:]...)
	}
	box.readUntilIndex = len(box.msgs) - 1
	return msgs
}

func (box *Box) GetAuth() *types.Auth {
	return box.auth
}

func (box *Box) append(msg *Msg) (int, error) {
	hash := msg.GetHash()
	_, exist := box.msgHashes[hash]
	if exist {
		return 0, code.AlreadyAppendedMsg
	}

	timestamp := msg.GetTimestamp()
	if box.latestTimestamp.Before(timestamp) {
		box.latestTimestamp = timestamp
	}

	box.msgs = append(box.msgs, msg)
	box.msgHashes[hash] = msg

	return len(box.msgs) - 1, nil
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
	box.personae[targetPersona.Address] = targetPersona
	return nil
}

func (box *Box) leave(targetPersona *types.Persona) error {
	oldPersona := box.getPersona(targetPersona.Address)
	if oldPersona == nil {
		return code.NonExistingPersonaInBox
	}
	delete(box.personae, targetPersona.Address)
	return nil
}

func (box *Box) propagate(auth *types.Auth, personae types.Personae) error {
	msg := NewMsg(box.myID, types.EmptyHash, &BodyUpdate{
		Auth:     auth,
		Personae: personae,
	})
	err := box.Publish(msg, TypeUpdate, true)
	if err != nil {
		return err
	}
	return nil
}
