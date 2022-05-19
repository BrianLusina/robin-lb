package lb

import (
	"log"
	"net/http"

	"github.com/brianlusina/robin-lb/app/internal/core/services"
	"github.com/brianlusina/robin-lb/app/tools"
)

// Handler is the load balancer handler
type Handler struct {
	http.HandlerFunc
	serverPool services.RobinService
}

// NewHandler creates a new Load Balancer Handler
func NewHandler(serverPool services.RobinService) *Handler {
	return &Handler{serverPool: serverPool}
}

// LoadBalance load balances incoming requests to the backends
func (lb *Handler) LoadBalance(w http.ResponseWriter, r *http.Request) {
	attempts := tools.GetAttemptsFromContext(r.Context())
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	lb.serverPool.ServeRequest(w, r)
}
