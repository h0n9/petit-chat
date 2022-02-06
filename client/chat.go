package client

import (
	"bufio"
	"fmt"
	"strings"
	"sync"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/control"
	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/msg"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type Chat struct {
	box *control.Box

	vault *types.Vault
	state *types.State
	store *msg.CapsuleStore

	wg            sync.WaitGroup
	chStopReceive chan bool
	chError       chan error
	chCapsuleSub  chan *msg.Capsule

	reader *bufio.Reader
}

func NewChat(box *control.Box, reader *bufio.Reader, nickname string, public bool) (*Chat, error) {
	privKey, err := crypto.GenPrivKey()
	if err != nil {
		return nil, err
	}
	secretKey, err := crypto.GenSecretKey()
	if err != nil {
		return nil, err
	}
	persona, err := types.NewPersona(nickname, nil, privKey.PubKey())
	if err != nil {
		return nil, err
	}
	chat := Chat{
		box: box,

		vault: types.NewVault(persona, privKey, secretKey),
		state: types.NewState(public),
		store: msg.NewCapsuleStore(),

		wg:            sync.WaitGroup{},
		chStopReceive: make(chan bool, 1),
		chError:       make(chan error, 1),
		chCapsuleSub:  nil,

		reader: reader,
	}
	err = chat.state.Join(persona)
	if err != nil {
		return nil, err
	}
	err = chat.state.Grant(persona, true, true, true)
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (c *Chat) setChCapsule(chCapsule chan *msg.Capsule) {
	c.chCapsuleSub = chCapsule
}

func (c *Chat) GetPersona(addr crypto.Addr) *types.Persona {
	return c.state.GetPersona(addr)
}

func (c *Chat) GetVault() *types.Vault {
	return c.vault
}

func (c *Chat) GetState() *types.State {
	return c.state
}

func (c *Chat) GetStore() *msg.CapsuleStore {
	return c.store
}

func (c *Chat) GetPeerID() types.ID {
	return c.box.GetHostID()
}

func (c *Chat) Close() {
	close(c.chStopReceive)
	close(c.chError)
}

func (c *Chat) Subscribe() error {
	c.setChCapsule(c.box.GetChCapsule())
	return c.box.Subscribe()
}

func (c *Chat) Stop() {
	c.setChCapsule(nil)
}

func (c *Chat) Send() {
	var stop bool = false
	for {
		fmt.Printf("> ")
		input, err := util.GetInput(c.reader, false, true)
		if err != nil {
			c.chError <- err
			continue
		}
		switch input {
		case "/exit":
			c.chStopReceive <- true
			stop = true
		case "/msgs":
			capsules := c.store.GetCapsules()
			for _, capsule := range capsules {
				if capsule.Encrypted {
					err = capsule.Decrypt(c.vault.GetSecretKey())
					if err != nil {
						c.chError <- err
						continue
					}
				}
				err = capsule.Check()
				if err != nil {
					c.chError <- err
					continue
				}
				m, err := capsule.Decapsulate()
				if err != nil {
					c.chError <- err
					continue
				}
				c.PrintMsg(m)
			}
			continue
		case "/peers":
			peers := c.box.GetPersonae()
			for _, peer := range peers {
				printPeer(peer)
			}
			continue
		case "/auth":
			auth := c.box.GetAuth()
			printAuth(auth)
			continue
		case "/grant":
			fmt.Printf("<address> <R|W|X>: ")
			input, err := util.GetInput(c.reader, false, false)
			if err != nil {
				c.chError <- err
				continue
			}

			// parse strings
			strs := strings.Split(input, " ")
			if len(strs) != 2 {
				continue
			}
			addr := crypto.Addr(strs[0])
			if len(addr) != crypto.AddrSize {
				c.chError <- code.ImproperAddress
				continue
			}
			r, w, x := parsePerm(strs[1])

			err = c.box.Grant(addr, r, w, x)
			if err != nil {
				c.chError <- err
				continue
			}
			continue
		case "/revoke":
			fmt.Printf("<address>: ")
			input, err := util.GetInput(c.reader, false, false)
			if err != nil {
				c.chError <- err
				continue
			}

			addr := crypto.Addr(input)
			if len(addr) != crypto.AddrSize {
				c.chError <- code.ImproperAddress
				continue
			}

			err = c.box.Revoke(addr)
			if err != nil {
				c.chError <- err
				continue
			}
			continue
		case "":
			continue
		}
		if stop {
			break
		}

		// CLI supports ONLY TypeText
		peerID := c.GetPeerID()
		clientAddr := c.vault.GetAddr()
		m := msg.NewMsgRaw(peerID, clientAddr, types.EmptyHash, []byte(input), nil)
		err = c.Publish(m, false)
		if err != nil {
			c.chError <- err
			continue
		}
	}
	c.wg.Done()
}

func (c *Chat) Receive() {
	var (
		stop    bool  = false
		err     error = nil
		capsule *msg.Capsule
	)
	for {
		select {
		case capsule = <-c.chCapsuleSub:
			m, err := c.Handler(capsule)
			if err != nil {
				c.chError <- err
				continue
			}
			c.PrintMsg(m)
		case err = <-c.chError:
			fmt.Println(err)
		case <-c.chStopReceive:
			stop = true
		}
		if stop {
			break
		}
	}
	c.wg.Done()
}

func (c *Chat) Publish(m *msg.Msg, encrypt bool) error {
	capsule, err := m.Encapsulate()
	if err != nil {
		return err
	}
	err = capsule.Sign(c.vault.GetPrivKey())
	if err != nil {
		return err
	}
	if encrypt {
		err = capsule.Encrypt(c.vault.GetSecretKey())
		if err != nil {
			return err
		}
	}
	err = c.box.Publish(capsule)
	if err != nil {
		return err
	}

	return nil
}

func (c *Chat) ReadMsg(m *msg.Msg, hash types.Hash) error {
	if m.GetType() <= msg.TypeMeta {
		return nil
	}
	vault := c.GetVault()
	if vault == nil {
		return code.ImproperVault
	}
	if m.GetClientAddr() == vault.GetAddr() {
		return nil
	}
	peerID := c.GetPeerID()
	clientAddr := vault.GetAddr()
	meta := types.NewMeta(false, true, false)
	msgMeta := msg.NewMsgMeta(peerID, clientAddr, types.EmptyHash, hash, meta)
	err := c.Publish(msgMeta, true)
	if err != nil {
		return err
	}
	return nil
}

func (c *Chat) PrintMsg(m *msg.Msg) {
	timestamp := m.GetTimestamp()
	addr := m.GetClientAddr()
	persona := c.GetPersona(addr)
	nickname := "somebody"
	if persona != nil {
		nickname = persona.GetNickname()
	}
	switch m.GetType() {
	case msg.TypeRaw:
		body := m.GetBody().(msg.BodyRaw)
		metas := m.GetMetas()
		fmt.Printf("[%s, %s] %s\n", timestamp, nickname, body.Data)
		for addr, meta := range metas {
			nickname = c.GetPersona(addr).Nickname
			fmt.Printf("  - %s %s\n", nickname, printMeta(meta))
		}
	case msg.TypeHelloSyn:
		fmt.Printf("[%s, %s] entered\n", timestamp, nickname)
	case msg.TypeHelloAck:
	case msg.TypeBye:
		fmt.Printf("[%s, %s] left\n", timestamp, nickname)
	case msg.TypeUpdate:
	case msg.TypeMeta:
		body := m.GetBody().(msg.BodyMeta)
		done := printMeta(body.Meta)
		fmt.Printf("[%s, %s] %s %x\n", timestamp, nickname, done, m.GetParentHash())
	default:
		fmt.Println("Unknown Type")
	}
}
