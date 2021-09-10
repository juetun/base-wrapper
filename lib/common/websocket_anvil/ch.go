package websocket_anvil

import (
	"sync"
)

type Chan struct {
	C      chan interface{}
	closed bool
	lock   sync.Mutex
}

func NewCh() *Chan {
	return &Chan{C: make(chan interface{})}
}

func (s *Chan) SafeClose() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.closed {
		close(s.C)
		s.closed = true
	}
}

func (s *Chan) IsClosed() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.closed
}

func (s *Chan) SafeSend(data interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.closed {
		s.C <- data
	}
}

