package cluster

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSignal(t *testing.T) {
	s := NewSignal(false)
	assert.False(t, s.IsSet())
	s.Set()
	assert.True(t, s.IsSet())
	s.Set()
	assert.True(t, s.IsSet())
}
