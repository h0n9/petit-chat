package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPerm(t *testing.T) {
	// case: all false
	p := NewPerm(false, false, false)
	assert.Equal(t, false, p.CanRead())
	assert.Equal(t, false, p.CanWrite())
	assert.Equal(t, false, p.CanExecute())

	// case: one true
	p = NewPerm(true, false, false)
	assert.Equal(t, true, p.CanRead())
	assert.Equal(t, false, p.CanWrite())
	assert.Equal(t, false, p.CanExecute())

	p = NewPerm(false, true, false)
	assert.Equal(t, false, p.CanRead())
	assert.Equal(t, true, p.CanWrite())
	assert.Equal(t, false, p.CanExecute())

	p = NewPerm(false, false, true)
	assert.Equal(t, false, p.CanRead())
	assert.Equal(t, false, p.CanWrite())
	assert.Equal(t, true, p.CanExecute())

	// case: two true
	p = NewPerm(true, true, false)
	assert.Equal(t, true, p.CanRead())
	assert.Equal(t, true, p.CanWrite())
	assert.Equal(t, false, p.CanExecute())

	p = NewPerm(true, false, true)
	assert.Equal(t, true, p.CanRead())
	assert.Equal(t, false, p.CanWrite())
	assert.Equal(t, true, p.CanExecute())

	p = NewPerm(false, true, true)
	assert.Equal(t, false, p.CanRead())
	assert.Equal(t, true, p.CanWrite())
	assert.Equal(t, true, p.CanExecute())

	// calse: all true
	p = NewPerm(true, true, true)
	assert.Equal(t, true, p.CanRead())
	assert.Equal(t, true, p.CanWrite())
	assert.Equal(t, true, p.CanExecute())
}
