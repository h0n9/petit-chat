package types

import (
	"encoding/json"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
)

type Auth struct {
	Public bool                 `json:"public"`
	Perms  map[crypto.Addr]Perm `json:"perms"`
}

func NewAuth(public bool, perms map[crypto.Addr]Perm) *Auth {
	return &Auth{
		Public: public,
		Perms:  perms,
	}
}

func (a *Auth) IsPublic() bool {
	return a.Public
}

func (a *Auth) getPerm(addr crypto.Addr) (Perm, error) {
	perm, exist := a.Perms[addr]
	if !exist {
		return 0, code.NonExistingPermission
	}
	return perm, nil
}

func (a *Auth) GetPerm(addr crypto.Addr) Perm {
	perm, err := a.getPerm(addr)
	if err != nil {
		return permNone
	}
	return perm
}

func (a *Auth) SetPerm(addr crypto.Addr, Perm Perm) error {
	// TODO: add constraints
	a.Perms[addr] = Perm
	return nil
}

func (a *Auth) SetPerms(Perms map[crypto.Addr]Perm) error {
	// TODO: add constraints
	a.Perms = Perms
	return nil
}

func (a *Auth) DeletePerm(addr crypto.Addr) error {
	_, err := a.getPerm(addr)
	if err != nil {
		return err
	}
	delete(a.Perms, addr)
	return nil
}

func (a *Auth) checkPerm(addr crypto.Addr, perm Perm) (bool, error) {
	// check id in perms first
	p, err := a.getPerm(addr)
	if err != nil {
		return false, err
	}
	return p.canDo(perm), nil
}

func (a *Auth) CanRead(addr crypto.Addr) bool {
	ok, err := a.checkPerm(addr, permRead)
	if err != nil {
		return false
	}
	return ok
}

func (a *Auth) CanWrite(addr crypto.Addr) bool {
	ok, err := a.checkPerm(addr, permWrite)
	if err != nil {
		return false
	}
	return ok
}

func (a *Auth) CanExecute(addr crypto.Addr) bool {
	ok, err := a.checkPerm(addr, permExecute)
	if err != nil {
		return false
	}
	return ok
}

func (a *Auth) Grant(persona *Persona, r, w, x bool) error {
	return a.SetPerm(persona.Address, NewPerm(r, w, x))
}

func (oldAuth *Auth) Copy() (*Auth, error) {
	var newAuth Auth
	data, err := json.Marshal(oldAuth)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &newAuth)
	if err != nil {
		return nil, err
	}
	return &newAuth, nil
}
