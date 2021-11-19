package types

import (
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

func (a *Auth) GetPerm(addr Addr) (Perm, error) {
	p, exist := a.Perms[addr]
	if !exist {
		return 0, code.NonExistingPermission
	}
	return p, nil
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

func (a *Auth) CheckMinPerm(addr Addr) (bool, error) {
	ok := false

	// check id in perms first
	p, err := a.GetPerm(addr)
	if err != nil {
		return ok, err
	}

	// check id's minimum permission
	switch a.Public {
	case true:
		ok = p.canDo(minPermPublic)
	case false:
		ok = p.canDo(minPermPrivate)
	}

	return ok, nil
}
