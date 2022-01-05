package msg

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
	"github.com/libp2p/go-libp2p-core/peer"
)

// Box refers to a chat room
type Box struct {
	ctx      context.Context
	sub      *types.Sub
	msgSubCh chan *Msg

	vault struct {
		hostID      types.ID
		hostPersona *types.Persona // TODO: move to client side permanently
		privKey     *crypto.PrivKey
		secretKey   *crypto.SecretKey
	}

	state struct {
		topic           *types.Topic
		personae        types.Personae
		auth            *types.Auth
		latestTimestamp time.Time
		readUntilIndex  int
	}

	store struct {
		msgs      []*Msg              // TODO: limit the size of msgs slice
		msgHashes map[types.Hash]*Msg // TODO: limit the size of msgHashes map
	}
}

func NewBox(ctx context.Context, topic *types.Topic, public bool,
	hostID types.ID, privKey *crypto.PrivKey, hostPersona *types.Persona) (*Box, error) {
	err := hostPersona.Check()
	if err != nil {
		return nil, err
	}
	secretKey, err := crypto.GenSecretKey()
	if err != nil {
		return nil, err
	}
	box := Box{
		ctx:      ctx,
		sub:      nil,
		msgSubCh: nil,

		vault: struct {
			hostID      peer.ID
			hostPersona *types.Persona
			privKey     *crypto.PrivKey
			secretKey   *crypto.SecretKey
		}{
			hostID:      hostID,
			hostPersona: hostPersona,
			privKey:     privKey,
			secretKey:   secretKey,
		},

		state: struct {
			topic           *types.Topic
			personae        types.Personae
			auth            *types.Auth
			latestTimestamp time.Time
			readUntilIndex  int
		}{
			topic:           topic,
			personae:        make(types.Personae),
			auth:            types.NewAuth(public, make(map[crypto.Addr]types.Perm)),
			latestTimestamp: time.Now(),
			readUntilIndex:  0,
		},

		store: struct {
			msgs      []*Msg
			msgHashes map[types.Hash]*Msg
		}{
			msgs:      make([]*Msg, 0),
			msgHashes: make(map[types.Hash]*Msg),
		},
	}
	err = box.join(hostPersona)
	if err != nil {
		return nil, err
	}
	err = grant(box.state.auth, hostPersona.Address, true, true, true)
	if err != nil {
		return nil, err
	}
	msg := NewMsgHelloSyn(&box, types.EmptyHash, hostPersona)
	err = box.Publish(msg, false)
	if err != nil {
		return nil, err
	}
	return &box, nil
}

func (box *Box) Encapsulate(msg *Msg, encrypt bool) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	if encrypt {
		data, err = box.vault.secretKey.Encrypt(data)
		if err != nil {
			return nil, err
		}
	}

	data, err = json.Marshal(&MsgCapsule{
		Encrypted: encrypt,
		Type:      msg.GetType(),
		Data:      data,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (box *Box) Decapsulate(data []byte) (*Msg, error) {
	msgCapsule := MsgCapsule{}
	err := json.Unmarshal(data, &msgCapsule)
	if err != nil {
		return nil, err
	}

	if msgCapsule.Encrypted {
		msgCapsule.Data, err = box.vault.secretKey.Decrypt(msgCapsule.Data)
		if err != nil {
			return nil, err
		}
	}

	msg := NewMsg(msgCapsule.Type.Base())
	if msg == nil {
		return nil, code.UnknownMsgType
	}

	err = json.Unmarshal(msgCapsule.Data, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (box *Box) Publish(msg *Msg, encrypt bool) error {
	err := box.Sign(msg)
	if err != nil {
		return err
	}
	data, err := box.Encapsulate(msg, encrypt)
	if err != nil {
		return err
	}
	err = box.state.topic.Publish(box.ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (box *Box) Subscribe(handler MsgHandler) error {
	if box.sub != nil {
		return code.AlreadySubscribingTopic
	}

	sub, err := box.state.topic.Subscribe()
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
		err = box.Verify(msg)
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
			err = box.state.topic.Close()
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}

	return nil
}

func Hash(base Base) types.Hash {
	data, err := json.Marshal(base)
	if err != nil {
		return types.EmptyHash
	}
	return util.ToSHA256(data)
}

func (box *Box) Sign(msg *Msg) error {
	data, err := json.Marshal(msg.Base)
	if err != nil {
		return err
	}
	sigBytes, err := box.vault.privKey.Sign(data)
	if err != nil {
		return err
	}
	msg.SetHash(util.ToSHA256(data))
	msg.SetSignature(Signature{
		SigBytes: sigBytes,
		PubKey:   box.vault.hostPersona.PubKey,
	})
	return nil
}

func (box *Box) Verify(msg *Msg) error {
	data, err := json.Marshal(msg.Base)
	if err != nil {
		return err
	}
	signature := msg.GetSignature()
	ok := signature.PubKey.Verify(data, signature.SigBytes)
	if !ok {
		return code.FailedToVerify
	}
	return nil
}

func (box *Box) GetPersonae() map[crypto.Addr]*types.Persona {
	return box.state.personae
}

func (box *Box) getPersona(cAddr crypto.Addr) *types.Persona {
	persona, exist := box.state.personae[cAddr]
	if !exist {
		return nil
	}
	return persona
}

func (box *Box) GetPersona(cAddr crypto.Addr) *types.Persona {
	return box.getPersona(cAddr)
}

func (box *Box) GetHostID() types.ID {
	return box.vault.hostID
}

func (box *Box) GetHostPersona() *types.Persona {
	return box.vault.hostPersona
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
	newAuth, err := box.state.auth.Copy()
	if err != nil {
		return err
	}
	err = grant(newAuth, addr, r, w, x)
	if err != nil {
		return err
	}

	err = box.propagate(newAuth, box.state.personae)
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
	newAuth, err := box.state.auth.Copy()
	if err != nil {
		return err
	}
	err = revoke(newAuth, addr)
	if err != nil {
		return err
	}

	err = box.propagate(newAuth, box.state.personae)
	if err != nil {
		return err
	}

	return nil
}

func (box *Box) Close() error {
	// Announe EOS to others (application layer)
	msg := NewMsgBye(box, types.EmptyHash, box.vault.hostPersona)
	return box.Publish(msg, true)
}

func (box *Box) Subscribing() bool {
	return box.sub != nil
}

func (box *Box) SetMsgSubCh(msgSubCh chan *Msg) {
	box.msgSubCh = msgSubCh
}

func (box *Box) GetSecretKey() *crypto.SecretKey {
	return box.vault.secretKey
}

func (box *Box) GetMsgs() []*Msg {
	return box.store.msgs
}

func (box *Box) GetMsg(mh types.Hash) *Msg {
	return box.store.msgHashes[mh]
}

func (box *Box) GetUnreadMsgs() []*Msg {
	msgs := []*Msg{}
	if box.state.readUntilIndex+1 < len(box.store.msgs) {
		msgs = append(msgs, box.store.msgs[box.state.readUntilIndex+1:]...)
	}
	box.state.readUntilIndex = len(box.store.msgs) - 1
	return msgs
}

func (box *Box) GetAuth() *types.Auth {
	return box.state.auth
}

func (box *Box) append(msg *Msg) (int, error) {
	hash := msg.GetHash()
	_, exist := box.store.msgHashes[hash]
	if exist {
		return 0, code.AlreadyAppendedMsg
	}

	timestamp := msg.GetTimestamp()
	if box.state.latestTimestamp.Before(timestamp) {
		box.state.latestTimestamp = timestamp
	}

	box.store.msgs = append(box.store.msgs, msg)
	box.store.msgHashes[hash] = msg

	return len(box.store.msgs) - 1, nil
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
	box.state.personae[targetPersona.Address] = targetPersona
	return nil
}

func (box *Box) leave(targetPersona *types.Persona) error {
	oldPersona := box.getPersona(targetPersona.Address)
	if oldPersona == nil {
		return code.NonExistingPersonaInBox
	}
	delete(box.state.personae, targetPersona.Address)
	return nil
}

func (box *Box) propagate(auth *types.Auth, personae types.Personae) error {
	msg := NewMsgUpdate(box, types.EmptyHash, auth, personae)
	err := box.Publish(msg, true)
	if err != nil {
		return err
	}
	return nil
}
