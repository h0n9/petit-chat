package types

import (
	"encoding/json"

	"github.com/h0n9/petit-chat/code"
)

const (
	minPermPublic  Perm = permNone
	minPermPrivate Perm = permRead
)

type Auth struct {
	Public bool          `json:"public"`
	Perms  map[Addr]Perm `json:"perms"`
}

func NewAuth(public bool, perms map[Addr]Perm) *Auth {
	return &Auth{
		Public: public,
		Perms:  perms,
	}
}

func (a *Auth) IsPublic() bool {
	return a.Public
}

func (a *Auth) getPerm(addr Addr) (Perm, error) {
	perm, exist := a.Perms[addr]
	if !exist {
		return 0, code.NonExistingPermission
	}
	return perm, nil
}

func (a *Auth) GetPerm(addr Addr) (Perm, error) {
	return a.getPerm(addr)
}

func (a *Auth) SetPerm(addr Addr, Perm Perm) error {
	// TODO: add constraints
	a.Perms[addr] = Perm
	return nil
}

func (a *Auth) SetPerms(Perms map[Addr]Perm) error {
	// TODO: add constraints
	a.Perms = Perms
	return nil
}

func (a *Auth) DeletePerm(addr Addr) error {
	_, err := a.getPerm(addr)
	if err != nil {
		return err
	}
	delete(a.Perms, addr)
	return nil
}

func (a *Auth) checkPerm(addr Addr, perm Perm) (bool, error) {
	// check id in perms first
	p, err := a.GetPerm(addr)
	if err != nil {
		return false, err
	}
	return p.canDo(perm), nil
}

func (a *Auth) CanRead(addr Addr) (bool, error) {
	return a.checkPerm(addr, permRead)
}

func (a *Auth) CanWrite(addr Addr) (bool, error) {
	return a.checkPerm(addr, permWrite)
}

func (a *Auth) CanExecute(addr Addr) (bool, error) {
	return a.checkPerm(addr, permExecute)
}

func (a *Auth) CheckMinPerm(addr Addr) (bool, error) {
	// check id's minimum permission
	if a.Public {
		return a.checkPerm(addr, minPermPublic)
	} else {
		return a.checkPerm(addr, minPermPrivate)
	}
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
