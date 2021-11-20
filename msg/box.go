package msg

import (
	"context"
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
	b := Box{
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

		personae:  make(map[crypto.Addr]*types.Persona),
		msgs:      make([]*Msg, 0),
		msgHashes: make(map[types.Hash]*Msg),
	}
	err = b.join(mp)
	if err != nil {
		return nil, err
	}
	msh := NewMsgStructHelloSyn(b.myPersona)
	data, err := msh.Encapsulate()
	if err != nil {
		return nil, err
	}
	err = b.Publish(MsgTypeHelloSyn, types.Hash{}, false, data)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (b *Box) Publish(t MsgType, parentMsgHash types.Hash, encrypt bool, data []byte) error {
	if len(data) == 0 {
		// this is not error
		return nil
	}
	if encrypt {
		encryptedData, err := b.secretKey.Encrypt(data)
		if err != nil {
			return err
		}
		data = encryptedData
	}
	msg, err := NewMsg(b.myID, b.myPersona.Address, t, parentMsgHash, encrypt, data)
	if err != nil {
		return err
	}
	err = msg.Sign(b.myPrivKey)
	if err != nil {
		return err
	}
	data, err = msg.Encapsulate()
	if err != nil {
		return err
	}
	err = b.topic.Publish(b.ctx, data)
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
		msg := new(Msg)
		err = msg.Decapsulate(data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = msg.Verify()
		if err != nil {
			fmt.Println(err)
			continue
		}
		eos, err := handler(b, msg, received.GetFrom())
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

func (b *Box) Close() error {
	// Announe EOS to others (application layer)
	msb := NewMsgStructBye(b.myPersona)
	data, err := msb.Encapsulate()
	if err != nil {
		return err
	}
	return b.Publish(MsgTypeBye, types.Hash{}, false, data)
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
