package robinsrv

import (
	"log"
	"net/http"
	"net/url"
	"sync/atomic"

	"github.com/brianlusina/robin-lb/app/internal/core/entities"
	"github.com/brianlusina/robin-lb/app/internal/core/services"
	"github.com/brianlusina/robin-lb/app/pkg"
)

type service struct {
	backends []*entities.Backend
	current  uint64
}

func New() *service {
	backends := []*entities.Backend{}
	return &service{
		backends: backends,
	}
}

// AddBackend adds a new backend to the server pool
func (s *service) AddBackend(serverURL *url.URL, alive bool, proxySvc services.ProxyService) {
	backend := entities.NewBackend(serverURL, alive, proxySvc)
	s.backends = append(s.backends, backend)
}

// NextIndex atomically increase the counter and return an index
// we increase the current value by one atomically and return the index by modding with the length of the slice.
// Which means the value always will be between 0 and length of the slice.
// In the end, we are interested in a particular index, not the total count.
func (s *service) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

// ServeRequest forwards the request to the next available backend
func (s *service) ServeRequest(w http.ResponseWriter, r *http.Request) {
	peer := s.GetNextActivePeer()
	if peer == nil {
		log.Println("No active peer found")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	peer.ReverseProxy.ServeHTTPRequest(w, r)
}

// GetNextPeer returns next active peer from the server pool
func (s *service) GetNextActivePeer() *entities.Backend {
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
func (s *service) HealthCheck() {
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

func (s *service) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	for _, b := range s.backends {
		if b.URL.String() == backendUrl.String() {
			b.SetAlive(alive)
			break
		}
	}
}
