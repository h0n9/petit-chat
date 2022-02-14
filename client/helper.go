package client

import (
	"fmt"
	"strings"

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

func printPersona(p *types.Persona) {
	fmt.Printf("[%s] %s\n", p.Address, p.Nickname)
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
