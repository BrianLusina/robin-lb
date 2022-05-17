package lb

import (
	"log"
	"net/http"

	"github.com/brianlusina/robin-lb/app/internal/core/services"
	"github.com/brianlusina/robin-lb/app/tools"
)

type lbHandler struct {
	http.HandlerFunc
	serverPool services.RobinService
}

func NewLbHandler(serverPool services.RobinService) *lbHandler {
	return &lbHandler{serverPool: serverPool}
}

// LoadBalance load balances incoming requests to the backends
func (lb *lbHandler) LoadBalance(w http.ResponseWriter, r *http.Request) {
	attempts := tools.GetAttemptsFromContext(r.Context())
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	lb.serverPool.ServeRequest(w, r)

	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}
