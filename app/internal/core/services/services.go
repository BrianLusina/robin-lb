package services

import (
	"net/http"
	"net/url"
)

type RobinService interface {
	ServeRequest(w http.ResponseWriter, r *http.Request)
	MarkBackendStatus(url *url.URL, status bool)
}

type ProxyService interface {
	ServeHTTPRequest(w http.ResponseWriter, r *http.Request)
}
