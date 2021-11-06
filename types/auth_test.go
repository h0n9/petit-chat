package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	a := NewAuth(false, nil)
	public := a.IsPublic()
	assert.False(t, public)
}
