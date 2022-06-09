// https://morioh.com/p/db6cac742e1c
package balancer

import (
	"sync"
	"sync/atomic"
)

type ServerPool struct {
	urls    []string
	current uint64
	mux     sync.RWMutex
}

// nextIndex atomically increase the counter and return an index
func (s *ServerPool) nextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.urls)))
}

func (s *ServerPool) GetNextAddr() string {
	s.mux.Lock()
	defer s.mux.Unlock()

	if len(s.urls) == 0 {
		return ""
	}

	next := s.nextIndex()
	return s.urls[next]
}

func (s *ServerPool) SetUrls(urls []string) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.urls = urls
	if s.current > uint64(len(s.urls)) {
		s.current = uint64(len(s.urls))
	}
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

func (n *NamesServerPool) GetNextAddr(name string) string {
	n.mux.Lock()
	defer n.mux.Unlock()

	pool, has := n.pools[name]
	if !has {
		return ""
		//n.pools[name] = &ServerPool{urls: make([]string, 0), current: 0}
	}

	return pool.GetNextAddr()
}

func (n *NamesServerPool) SetUrls(name string, urls []string) {
	n.mux.Lock()
	defer n.mux.Unlock()

	_, has := n.pools[name]
	if !has {
		n.pools[name] = &ServerPool{urls: make([]string, 0), current: 0}
	}

	n.pools[name].SetUrls(urls)
}
