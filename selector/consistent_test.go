package selector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	ch := NewWith(20)
	ch.Add(&addrValue{"tcp", "127.0.0.1:11211"})

	assert.Equal(t, 20, ch.circle.Size())
}

func TestRemove(t *testing.T) {
	ch := NewWith(20)

	ch.Add(&addrValue{"tcp", "127.0.0.1:11211"})
	assert.Equal(t, 20, ch.circle.Size())

	ch.Remove(&addrValue{"tcp", "127.0.0.1:11211"})
	assert.Equal(t, 0, ch.circle.Size())
}

func TestPickServer1Server(t *testing.T) {
	ch := NewWith(20)
	server := &addrValue{"tcp", "127.0.0.1:11211"}
	ch.Add(server)

	picked, _ := ch.PickForKey("any key")

	assert.EqualValues(t, server, picked)
}
