package client

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/msg"
	"github.com/h0n9/petit-chat/types"
)

type MsgHandler func(m *msg.Msg, hostID types.ID) (bool, error)

func (c *Chat) Handler(capsule *msg.Capsule) (*msg.Msg, error) {
	if capsule.Encrypted {
		err := capsule.Decrypt(c.vault.GetSecretKey())
		if err != nil {
			return nil, err
		}
	}
	err := capsule.Check()
	if err != nil {
		return nil, err
	}
	m, err := capsule.Decapsulate()
	if err != nil {
		return nil, err
	}

	hash := capsule.GetHash()
	// msg handling flow:
	//   check -> append -> execute -> (received)

	// check msg.Body
	err = m.Check(hash, c)
	if err != nil && err != code.SelfMsg {
		return nil, err
	}

	// execute msg.Body
	err = m.Execute(hash, c)
	if err != nil {
		return nil, err
	}

	index, err := c.store.Append(capsule)
	if err != nil {
		return nil, err
	}
	c.state.SetReadUntilIndex(index)

	if m.GetType() <= msg.TypeMeta || m.GetClientAddr() == c.vault.GetAddr() {
		return m, nil
	}

	// meta := types.NewMeta(true, canRead, false)
	// msgMeta := NewMsgMeta(box, types.EmptyHash, msg.GetHash(), meta)
	// err = box.Publish(msgMeta, true)
	// if err != nil {
	// 	return eos, err
	// }
	return m, nil
}
