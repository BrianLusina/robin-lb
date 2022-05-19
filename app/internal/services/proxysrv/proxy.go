package proxysrv

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/brianlusina/robin-lb/app/internal/core/services"
	"github.com/brianlusina/robin-lb/app/tools"
)

type service struct {
	*httputil.ReverseProxy
	serverURL  *url.URL
	serverPool services.RobinService
	handler    http.HandlerFunc
}

// New returns a new Proxy service
func New(serverURL *url.URL, handler http.HandlerFunc, serverPool services.RobinService) *service {
	return &service{
		ReverseProxy: httputil.NewSingleHostReverseProxy(serverURL),
		serverURL:    serverURL,
		serverPool:   serverPool,
	}
}

// ServeHTTPRequest is the reverse proxy handler
func (p *service) ServeHTTPRequest(w http.ResponseWriter, r *http.Request) {
	p.ServeHTTP(w, r)
}

// AddErroHandler adds an error handler to the revers proxy
func (p *service) AddErrorHandler() {
	p.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		log.Printf("[%s] %s\n", p.serverURL.Host, e.Error())
		retries := tools.GetRetryFromContext(request.Context())

		if retries < 3 {
			for range time.After(10 * time.Millisecond) {
				ctx := context.WithValue(request.Context(), tools.RetryKey, retries+1)
				p.ServeHTTP(writer, request.WithContext(ctx))
			}
			return
		}

		// after 3 retries, mark this backend as down
		p.serverPool.MarkBackendStatus(p.serverURL, false)

		// if the same request routing for few attempts with different backends, increase the count
		attempts := tools.GetAttemptsFromContext(request.Context())

		log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
		ctx := context.WithValue(request.Context(), tools.AttemptsKey, attempts+1)

		p.handler(writer, request.WithContext(ctx))
	}
}
