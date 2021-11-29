package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPerm(t *testing.T) {
	// case: all false
	p := NewPerm(false, false, false)
	assert.Equal(t, false, p.canDo(permRead))
	assert.Equal(t, false, p.canDo(permWrite))
	assert.Equal(t, false, p.canDo(permExecute))

	// case: one true
	p = NewPerm(true, false, false)
	assert.Equal(t, true, p.canDo(permRead))
	assert.Equal(t, false, p.canDo(permWrite))
	assert.Equal(t, false, p.canDo(permExecute))

	p = NewPerm(false, true, false)
	assert.Equal(t, false, p.canDo(permRead))
	assert.Equal(t, true, p.canDo(permWrite))
	assert.Equal(t, false, p.canDo(permExecute))

	p = NewPerm(false, false, true)
	assert.Equal(t, false, p.canDo(permRead))
	assert.Equal(t, false, p.canDo(permWrite))
	assert.Equal(t, true, p.canDo(permExecute))

	// case: two true
	p = NewPerm(true, true, false)
	assert.Equal(t, true, p.canDo(permRead))
	assert.Equal(t, true, p.canDo(permWrite))
	assert.Equal(t, false, p.canDo(permExecute))

	p = NewPerm(true, false, true)
	assert.Equal(t, true, p.canDo(permRead))
	assert.Equal(t, false, p.canDo(permWrite))
	assert.Equal(t, true, p.canDo(permExecute))

	p = NewPerm(false, true, true)
	assert.Equal(t, false, p.canDo(permRead))
	assert.Equal(t, true, p.canDo(permWrite))
	assert.Equal(t, true, p.canDo(permExecute))

	// calse: all true
	p = NewPerm(true, true, true)
	assert.Equal(t, true, p.canDo(permRead))
	assert.Equal(t, true, p.canDo(permWrite))
	assert.Equal(t, true, p.canDo(permExecute))
}
