package types

type Perm uint8

const (
	permNone    Perm = 0 // 0000 0000
	permRead    Perm = 1 // 0000 0001
	permWrite   Perm = 2 // 0000 0010
	permExecute Perm = 4 // 0000 0100
)

func NewPerm(read, write, execute bool) Perm {
	var p Perm
	if read {
		p |= permRead
	}
	if write {
		p |= permWrite
	}
	if execute {
		p |= permExecute
	}

	return p
}

func (p Perm) canDo(pc Perm) bool {
	return p&pc == pc
}

func (p Perm) CanRead() bool {
	return p.canDo(permRead)
}

func (p Perm) CanWrite() bool {
	return p.canDo(permWrite)
}

func (p Perm) CanExecute() bool {
	return p.canDo(permExecute)
}
