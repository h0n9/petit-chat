package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMeta(t *testing.T) {
	meta := NewMeta(false, false, false)
	assert.Equal(t, false, meta.Received())
	assert.Equal(t, false, meta.Read())
	assert.Equal(t, false, meta.Typing())

	meta = NewMeta(true, false, false)
	assert.Equal(t, true, meta.Received())
	assert.Equal(t, false, meta.Read())
	assert.Equal(t, false, meta.Typing())

	meta = NewMeta(false, true, false)
	assert.Equal(t, false, meta.Received())
	assert.Equal(t, true, meta.Read())
	assert.Equal(t, false, meta.Typing())

	meta = NewMeta(false, false, true)
	assert.Equal(t, false, meta.Received())
	assert.Equal(t, false, meta.Read())
	assert.Equal(t, true, meta.Typing())

	meta = NewMeta(true, true, false)
	assert.Equal(t, true, meta.Received())
	assert.Equal(t, true, meta.Read())
	assert.Equal(t, false, meta.Typing())

	meta = NewMeta(false, true, true)
	assert.Equal(t, false, meta.Received())
	assert.Equal(t, true, meta.Read())
	assert.Equal(t, true, meta.Typing())

	meta = NewMeta(true, true, true)
	assert.Equal(t, true, meta.Received())
	assert.Equal(t, true, meta.Read())
	assert.Equal(t, true, meta.Typing())
}
