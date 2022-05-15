package entities

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func NewBackend(url *url.URL, alive bool, proxy *httputil.ReverseProxy) *Backend {
	return &Backend{
		URL:          url,
		Alive:        alive,
		ReverseProxy: proxy,
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
