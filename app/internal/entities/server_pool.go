package entities

import (
	"log"
	"net/url"
	"sync/atomic"

	"github.com/brianlusina/robin-lb/app/pkg"
)

type ServerPool struct {
	backends []*Backend
	current  uint64
}

// AddBackend to the server pool
func (s *ServerPool) AddBackend(backend *Backend) {
	s.backends = append(s.backends, backend)
}

// NextIndex atomically increase the counter and return an index
// we increase the current value by one atomically and return the index by modding with the length of the slice.
// Which means the value always will be between 0 and length of the slice.
// In the end, we are interested in a particular index, not the total count.
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

// GetNextPeer returns next active peer from the server pool
func (s *ServerPool) GetNextActivePeer() *Backend {
	// loop entire backends to find out an alive backend
	next := s.NextIndex()

	// start from next and move a full cycle
	l := len(s.backends) + next

	for i := next; i < l; i++ {
		// take an index by modding with length
		idx := i % len(s.backends)

		// if we have an alive backend, use it and store if its not the original one
		if s.backends[idx].IsAlive() {
			if i != next {
				// mark the current one
				atomic.StoreUint64(&s.current, uint64(idx))
			}
			return s.backends[idx]
		}
	}
	return nil
}

// HealthCheck pings the backends and updates the status
func (s *ServerPool) HealthCheck() {
	for _, b := range s.backends {
		status := "up"
		alive := pkg.IsBackendAlive(b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}

		log.Printf("%s [%s]\n", b.URL, status)
	}
}

func (s *ServerPool) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	for _, b := range s.backends {
		if b.URL.String() == backendUrl.String() {
			b.SetAlive(alive)
			break
		}
	}
}
