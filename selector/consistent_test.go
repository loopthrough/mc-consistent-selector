package selector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	ch := NewWith(20)
	ch.Add("test server")

	assert.Equal(t, 20, ch.circle.Size())
}

func TestRemove(t *testing.T) {
	ch := NewWith(20)

	ch.Add("test server")
	assert.Equal(t, 20, ch.circle.Size())

	ch.Remove("test server")
	assert.Equal(t, 0, ch.circle.Size())
}
