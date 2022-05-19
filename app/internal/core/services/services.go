package services

import (
	"net/http"
	"net/url"
)

// RobinService is the interface for the backend server pool
type RobinService interface {
	ServeRequest(w http.ResponseWriter, r *http.Request)
	MarkBackendStatus(url *url.URL, status bool)
}

// ProxyService is the interface for the reverse proxy
type ProxyService interface {
	ServeHTTPRequest(w http.ResponseWriter, r *http.Request)
}
