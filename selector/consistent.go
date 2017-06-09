package selector

import (
	"errors"
	"github.com/emirpasic/gods/maps/treemap"
	"hash/crc32"
	"net"
	"strconv"
	"sync"
)

var (
	ErrNoNodes = errors.New("No nodes added.")
)

type ConsistentHash struct {
	mu            sync.RWMutex
	circle        *treemap.Map
	pointsPerNode int
	nodes         []string
}

func New() *ConsistentHash {
	return &ConsistentHash{
		circle:        treemap.NewWithIntComparator(),
		pointsPerNode: 100, // 100 for decent distribution
	}
}

func NewWith(pointsPerNode int) *ConsistentHash {
	return &ConsistentHash{
		circle:        treemap.NewWithIntComparator(),
		pointsPerNode: pointsPerNode,
	}
}

// hash computes hash of the provided value
func hash(value string) uint32 {
	bbuf := keyBufPool.Get().(*[]byte)
	n := copy(*bbuf, value)
	return crc32.ChecksumIEEE((*bbuf)[:n])
}

// indexedKeyHash computes hash of the key with provided index
// in the form hash(key_index)
func indexedKeyHash(key string, index int) uint32 {
	indexedKey := key + "_" + strconv.Itoa(index)
	return hash(indexedKey)
}

// Add inserts points in the circle for the provided server
func (ch *ConsistentHash) Add(server net.Addr) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	for i := 0; i < ch.pointsPerNode; i++ {
		ch.circle.Put(
			int(indexedKeyHash(server.String(), i)), // cast for TreeMap comparisons
			server)
	}
}

// Remove removes points in the circle for the provided server
func (ch *ConsistentHash) Remove(server net.Addr) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	for i := 0; i < ch.pointsPerNode; i++ {
		ch.circle.Remove(
			int(indexedKeyHash(server.String(), i)))
	}
}

// PickForKey returns the value for key that ensures consistent distribution
func (ch *ConsistentHash) PickForKey(key string) (net.Addr, error) {
	if ch.circle.Size() == 0 {
		return nil, ErrNoNodes
	}

	h := int(hash(key))
	if server, found := ch.circle.Get(h); found {
		return server.(net.Addr), nil
	}
	upperMap := ch.circle.Select(func(key interface{}, value interface{}) bool {
		return key.(int) > int(h)
	})

	if upperMap.Size() == 0 {
		keys := ch.circle.Keys()
		server, _ := ch.circle.Get(keys[0])
		return server.(net.Addr), nil
	}

	keys := upperMap.Keys()
	server, _ := upperMap.Get(keys[0])
	return server.(net.Addr), nil
}

// keyBufPool is providing storage for
// preallocated byte slices so that they
// get reused when needed for hashing.
var keyBufPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 256)
		return &b
	},
}
