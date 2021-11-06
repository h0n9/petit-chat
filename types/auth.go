package types

import "github.com/h0n9/petit-chat/code"

const (
	PUBLIC_MIN_PERM   Perm = permNone
	PRIAVATE_MIN_PERM Perm = permReadable
)

type Auth struct {
	Public bool        `json:"public"`
	Perms  map[ID]Perm `json:"perms"`
}

func NewAuth(public bool, perms map[ID]Perm) *Auth {
	return &Auth{
		Public: public,
		Perms:  perms,
	}
}

func (a *Auth) IsPublic() bool {
	return a.Public
}

func (a *Auth) GetPerm(id ID) (Perm, error) {
	p, exist := a.Perms[id]
	if !exist {
		return 0, code.NonExistingPermission
	}
	return p, nil
}

func (a *Auth) SetPerm(id ID, Perm Perm) error {
	// TODO: add constraints
	a.Perms[id] = Perm
	return nil
}

func (a *Auth) SetPerms(Perms map[ID]Perm) error {
	// TODO: add constraints
	a.Perms = Perms
	return nil
}
