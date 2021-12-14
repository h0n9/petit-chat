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

func NewBox(ctx context.Context, tp *types.Topic, pub bool, mi types.ID, mpk *crypto.PrivKey, mp *types.Persona) (*Box, error) {
	err := mp.Check()
	if err != nil {
		return nil, err
	}
	secretKey, err := crypto.GenSecretKey()
	if err != nil {
		return nil, err
	}
	box := Box{
		ctx:       ctx,
		topic:     tp,
		sub:       nil,
		auth:      types.NewAuth(pub, make(map[crypto.Addr]types.Perm)),
		secretKey: secretKey,

		myID:            mi,
		myPrivKey:       mpk,
		myPersona:       mp,
		msgSubCh:        nil,
		latestTimestamp: time.Now(),
		readUntilIndex:  0,

		personae:  make(types.Personae),
		msgs:      make([]*Msg, 0),
		msgHashes: make(map[types.Hash]*Msg),
	}
	err = box.join(mp)
	if err != nil {
		return nil, err
	}
	err = grant(box.auth, box.myPersona.Address, true, true, true)
	if err != nil {
		return nil, err
	}
	msg, err := NewMsg(mi, mp.Address, types.EmptyHash, &BodyHelloSyn{
		Persona: mp,
	})
	if err != nil {
		return nil, err
	}
	err = box.Publish(msg, TypeHelloSyn, false)
	if err != nil {
		return nil, err
	}
	return &box, nil
}

func (b *Box) Encapsulate(msg *Msg, msgType Type, encrypt bool) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	if encrypt {
		data, err = b.secretKey.Encrypt(data)
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

func (b *Box) Decapsulate(data []byte) (*Msg, error) {
	msgCapsule := MsgCapsule{}
	msg := Msg{}

	err := json.Unmarshal(data, &msgCapsule)
	if err != nil {
		return nil, err
	}

	if msgCapsule.Encrypted {
		msgCapsule.Data, err = b.secretKey.Decrypt(msgCapsule.Data)
		if err != nil {
			return nil, err
		}
	}

	msg.Body = msgCapsule.Type.Body()

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

func (b *Box) Subscribe(handler MsgHandler) error {
	if b.sub != nil {
		return code.AlreadySubscribingTopic
	}

	sub, err := b.topic.Subscribe()
	if err != nil {
		return err
	}
	b.sub = sub

	for {
		received, err := sub.Next(b.ctx)
		if err != nil {
			// TODO: replace fmt.Println() to logger.Println()
			fmt.Println(err)
			continue
		}
		data := received.GetData()
		msg, err := b.Decapsulate(data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = msg.Verify()
		if err != nil {
			fmt.Println(err)
			continue
		}
		eos, err := handler(b, msg)
		if err != nil {
			// TODO: replace fmt.Println() to logger.Println()
			fmt.Println(err)
			continue
		}

		// eos shoud be the only way to break for loop
		if eos {
			b.sub.Cancel()
			err = b.topic.Close()
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}

	return nil
}

func (b *Box) GetPersonae() map[crypto.Addr]*types.Persona {
	return b.personae
}

func (b *Box) getPersona(cAddr crypto.Addr) *types.Persona {
	persona, exist := b.personae[cAddr]
	if !exist {
		return nil
	}
	return persona
}

func (b *Box) GetPersona(cAddr crypto.Addr) *types.Persona {
	return b.getPersona(cAddr)
}

func (b *Box) GetMyID() types.ID {
	return b.myID
}

func (b *Box) GetMyPersona() *types.Persona {
	return b.myPersona
}

func grant(auth *types.Auth, addr crypto.Addr, r, w, x bool) error {
	perm := types.NewPerm(r, w, x)
	err := auth.SetPerm(addr, perm)
	if err != nil {
		return err
	}
	return nil
}

func (b *Box) Grant(addr crypto.Addr, r, w, x bool) error {
	newAuth, err := b.auth.Copy()
	if err != nil {
		return err
	}
	err = grant(newAuth, addr, r, w, x)
	if err != nil {
		return err
	}

	err = b.propagate(newAuth, b.personae)
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

func (b *Box) Revoke(addr crypto.Addr) error {
	newAuth, err := b.auth.Copy()
	if err != nil {
		return err
	}
	err = revoke(newAuth, addr)
	if err != nil {
		return err
	}

	err = b.propagate(newAuth, b.personae)
	if err != nil {
		return err
	}

	return nil
}

func (b *Box) Close() error {
	// Announe EOS to others (application layer)
	msg, err := NewMsg(b.myID, b.myPersona.Address, types.EmptyHash, &BodyBye{
		Persona: b.myPersona,
	})
	if err != nil {
		return err
	}
	return b.Publish(msg, TypeBye, true)
}

func (b *Box) Subscribing() bool {
	return b.sub != nil
}

func (b *Box) SetMsgSubCh(msgSubCh chan *Msg) {
	b.msgSubCh = msgSubCh
}

func (b *Box) GetSecretKey() *crypto.SecretKey {
	return b.secretKey
}

func (b *Box) GetMsgs() []*Msg {
	return b.msgs
}

func (b *Box) GetMsg(mh types.Hash) *Msg {
	return b.msgHashes[mh]
}

func (b *Box) GetUnreadMsgs() []*Msg {
	msgs := []*Msg{}
	if b.readUntilIndex+1 < len(b.msgs) {
		msgs = append(msgs, b.msgs[b.readUntilIndex+1:]...)
	}
	b.readUntilIndex = len(b.msgs) - 1
	return msgs
}

func (b *Box) GetAuth() *types.Auth {
	return b.auth
}

func (b *Box) append(msg *Msg) (int, error) {
	hash := msg.GetHash()
	_, exist := b.msgHashes[hash]
	if exist {
		return 0, code.AlreadyAppendedMsg
	}

	timestamp := msg.GetTimestamp()
	if b.latestTimestamp.Before(timestamp) {
		b.latestTimestamp = timestamp
	}

	b.msgs = append(b.msgs, msg)
	b.msgHashes[hash] = msg

	return len(b.msgs) - 1, nil
}

func (b *Box) join(targetPersona *types.Persona) error {
	oldPersona := b.getPersona(targetPersona.Address)
	if oldPersona != nil {
		return nil // ignore even if existing
		// return code.ExistingPersonaInBox
	}
	err := targetPersona.Check()
	if err != nil {
		return err
	}
	b.personae[targetPersona.Address] = targetPersona
	return nil
}

func (b *Box) leave(targetPersona *types.Persona) error {
	oldPersona := b.getPersona(targetPersona.Address)
	if oldPersona == nil {
		return code.NonExistingPersonaInBox
	}
	delete(b.personae, targetPersona.Address)
	return nil
}

func (b *Box) propagate(auth *types.Auth, personae types.Personae) error {
	msg, err := NewMsg(b.myID, b.myPersona.Address, types.EmptyHash, &BodyUpdate{
		Auth:     auth,
		Personae: personae,
	})
	if err != nil {
		return err
	}
	err = b.Publish(msg, TypeUpdate, true)
	if err != nil {
		return err
	}
	return nil
}
