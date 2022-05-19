package entities

import (
	"net/url"
	"sync"

	"github.com/brianlusina/robin-lb/app/internal/core/services"
)

// Backend is an actual backend added to the Server pool where requests are proxied from the load balancer
type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy services.ProxyService
}

// NewBackend creates a new backend
func NewBackend(url *url.URL, alive bool, reverseProxy services.ProxyService) *Backend {
	return &Backend{
		URL:          url,
		Alive:        alive,
		ReverseProxy: reverseProxy,
	}
}

// SetAlive sets the backend's alive status
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Alive = alive
}

// IsAlive returns the backend's alive status
func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Alive
}
