package selector

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net"
	"strconv"
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

func TestDistributeAccrossServers(t *testing.T) {
	server1 := &addrValue{"tcp", "127.0.0.1:11211"}
	server2 := &addrValue{"tcp", "127.0.0.2:11211"}
	server3 := &addrValue{"tcp", "127.0.0.3:11211"}
	ch := New()
	ch.Add(server1)
	ch.Add(server2)
	ch.Add(server3)

	picks := map[net.Addr]int{
		server1: 0,
		server2: 0,
		server3: 0,
	}

	for i := 0; i < 1000; i++ {
		picked, _ := ch.PickForKey(strconv.Itoa(rand.Int()))
		picks[picked] += 1
	}

	server1Picks, _ := picks[server1]
	server2Picks, _ := picks[server2]
	server3Picks, _ := picks[server2]
	assert.InDelta(t, 300, server1Picks, 100)
	assert.InDelta(t, 300, server2Picks, 100)
	assert.InDelta(t, 300, server3Picks, 100)
}
