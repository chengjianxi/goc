// https://morioh.com/p/db6cac742e1c
package balancer

import (
	"sync"
	"sync/atomic"
)

type ServerPool struct {
	addrs   []string
	current uint64
	mux     sync.RWMutex
}

// nextIndex atomically increase the counter and return an index
func (s *ServerPool) nextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.addrs)))
}

func (s *ServerPool) GetNextAddr() string {
	s.mux.Lock()
	defer s.mux.Unlock()

	if len(s.addrs) == 0 {
		return ""
	}

	next := s.nextIndex()
	return s.addrs[next]
}

func (s *ServerPool) SetAddrs(addrs []string) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.addrs = addrs
}

type NamesServerPool struct {
	pools map[string]*ServerPool
	mux   sync.RWMutex
}

func NewNamesServerPool() *NamesServerPool {
	return &NamesServerPool{
		pools: make(map[string]*ServerPool),
	}
}

func (n *NamesServerPool) GetServerAddr(name string) string {
	n.mux.Lock()
	defer n.mux.Unlock()

	pool, has := n.pools[name]
	if !has {
		return ""
	}

	return pool.GetNextAddr()
}

func (n *NamesServerPool) SetServerAddrs(name string, addrs []string) {
	n.mux.Lock()
	defer n.mux.Unlock()

	_, has := n.pools[name]
	if !has {
		n.pools[name] = &ServerPool{addrs: make([]string, 0), current: 0}
	}

	n.pools[name].SetAddrs(addrs)
}
