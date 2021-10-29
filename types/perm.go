package types

type Perm uint8

const (
	permReadable   Perm = 1 // 0000 0001
	permWritable   Perm = 2 // 0000 0010
	permExecutable Perm = 4 // 0000 0100
)

func NewPerm(read, write, execute bool) Perm {
	var p Perm
	if read {
		p |= permReadable
	}
	if write {
		p |= permWritable
	}
	if execute {
		p |= permExecutable
	}

	return p
}

func (p Perm) CanRead() bool {
	return p&permReadable == permReadable
}

func (p Perm) CanWrite() bool {
	return p&permWritable == permWritable
}

func (p Perm) CanExecute() bool {
	return p&permExecutable == permExecutable
}
