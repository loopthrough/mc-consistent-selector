package selector

import (
	"net"
	"strings"
	"sync"
)

type ServerList struct {
	mu             sync.RWMutex
	servers        []net.Addr
	consistentHash *ConsistentHash
}

func NewServerList() *ServerList {
	return &ServerList{
		consistentHash: NewConsistentHash(),
	}
}

// addrVal is static storage for values provided as net.Addr.
// It implements the same interface as net.Addr
type addrValue struct {
	network, addrString string
}

func newAddrValue(a net.Addr) net.Addr {
	return &addrValue{
		network:    a.Network(),
		addrString: a.String(),
	}
}

func (v *addrValue) Network() string { return v.network }
func (v *addrValue) String() string  { return v.addrString }

// Add adds new server if there is no error.
func (sl *ServerList) Add(server string) error {
	var serverAddress net.Addr
	if strings.Contains(server, "/") {
		addr, err := net.ResolveUnixAddr("unix", server)
		if err != nil {
			return err
		}

		serverAddress = newAddrValue(addr)
	} else {
		tcpaddr, err := net.ResolveTCPAddr("tcp", server)
		if err != nil {
			return err
		}
		serverAddress = newAddrValue(tcpaddr)
	}

	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.servers = append(sl.servers, serverAddress)
	sl.consistentHash.Add(serverAddress)
	return nil
}

// Clear resets all the data structures to zero values
func (sl *ServerList) Clear() {
	newSl := new(ServerList)
	sl.servers = newSl.servers
	sl.consistentHash = newSl.consistentHash
}

// Each iterates over each server calling the given function
func (sl *ServerList) Each(f func(net.Addr) error) error {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	for _, a := range sl.servers {
		if err := f(a); nil != err {
			return err
		}
	}
	return nil
}

func (sl *ServerList) PickServer(key string) (net.Addr, error) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	return sl.consistentHash.PickForKey(key)
}
