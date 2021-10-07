package types

type Perm struct {
	Read    bool `json:"read"`
	Write   bool `json:"write"`
	Execute bool `json:"execute"`
}

func NewPerm(Read, Write, Execute bool) *Perm {
	return &Perm{
		Read:    Read,
		Write:   Write,
		Execute: Execute,
	}
}

type Auth struct {
	IsPublic bool         `json:"is_public"`
	Perms    map[ID]*Perm `json:"perms"`
}

func NewAuth(isPublic bool, perms map[ID]*Perm) *Auth {
	return &Auth{
		IsPublic: isPublic,
		Perms:    perms,
	}
}

func (a *Auth) SetPerm(id ID, perm *Perm) error {
	// TODO: add constraints
	a.Perms[id] = perm
	return nil
}

func (a *Auth) SetPerms(perms map[ID]*Perm) error {
	// TODO: add constraints
	a.Perms = perms
	return nil
}

func (a *Auth) CanRead(id ID) bool {
	canRead := false
	if perm := a.getPerm(id); perm != nil {
		canRead = perm.Read
	}
	return canRead
}

func (a *Auth) CanWrite(id ID) bool {
	canWrite := false
	if perm := a.getPerm(id); perm != nil {
		canWrite = perm.Write
	}
	return canWrite
}

func (a *Auth) CanExecute(id ID) bool {
	canExecute := false
	if perm := a.getPerm(id); perm != nil {
		canExecute = perm.Execute
	}
	return canExecute
}

func (a *Auth) getPerm(id ID) *Perm {
	return a.Perms[id]
}
