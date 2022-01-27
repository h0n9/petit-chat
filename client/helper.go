package client

import (
	"fmt"
	"strings"

	"github.com/h0n9/petit-chat/msg"
	"github.com/h0n9/petit-chat/types"
)

func printAuth(a *types.Auth) {
	p := "private"
	if a.IsPublic() {
		p = "public"
	}
	str := fmt.Sprintf("Auth: %s\n", p)
	if len(a.Perms) > 0 {
		str += "Perms:\n"
	}
	for addr := range a.Perms {
		str += fmt.Sprintf("[%s] ", addr)
		if a.CanRead(addr) {
			str += "R"
		}
		if a.CanWrite(addr) {
			str += "W"
		}
		if a.CanExecute(addr) {
			str += "X"
		}
		str += "\n"
	}
	fmt.Printf("%s", str)
}

func printPeer(p *types.Persona) {
	fmt.Printf("[%s] %s\n", p.Address, p.Nickname)
}

func readMsg(b *msg.Box, m *msg.Msg) error {
	printMsg(b, m)
	if m.GetType() <= msg.TypeMeta {
		return nil
	}
	if m.GetClientAddr() == b.GetHostPersona().Address {
		return nil
	}
	meta := types.NewMeta(false, true, false)
	msgMeta := msg.NewMsgMeta(b, types.EmptyHash, m.GetHash(), meta)
	err := b.Publish(msgMeta, true)
	if err != nil {
		return err
	}
	return nil
}

func printMsg(box *msg.Box, m *msg.Msg) {
	timestamp := m.GetTimestamp()
	addr := m.GetSignature().PubKey.Address()
	persona := box.GetPersona(addr)
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
			nickname = box.GetPersona(addr).Nickname
			fmt.Printf("  - %s %s\n", nickname, printMeta(meta))
		}
	case msg.TypeHelloSyn:
		fmt.Printf("[%s, %s] entered\n", timestamp, nickname)
	case msg.TypeHelloAck:
	case msg.TypeBye:
		fmt.Printf("[%s, %s] left\n", timestamp, nickname)
	case msg.TypeUpdate:
	case msg.TypeMeta:
		// body := m.GetBody().(msg.BodyMeta)
		// done := printMeta(body.Meta)
		// fmt.Printf("[%s, %s] %s %x\n", timestamp, nickname, done, m.GetParentHash())
	default:
		fmt.Println("Unknown Type")
	}
}

func printMeta(meta types.Meta) string {
	str := ""
	if meta.Received() {
		str += "received,"
	}
	if meta.Read() {
		str += "read,"
	}
	if meta.Typing() {
		str += "typing,"
	}
	return str
}

func parsePerm(permStr string) (bool, bool, bool) {
	r, w, x := false, false, false
	permStr = strings.ToUpper(permStr)
	if strings.Contains(permStr, "R") {
		r = true
	}
	if strings.Contains(permStr, "W") {
		w = true
	}
	if strings.Contains(permStr, "X") {
		x = true
	}
	return r, w, x
}
