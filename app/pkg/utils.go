package pkg

import (
	"log"
	"net"
	"net/url"
	"time"
)

// IsBackendAlive checks whether a backend is Alive by establishing a TCP connection
func IsBackendAlive(url *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", url.Host, timeout)

	if err != nil {
		log.Printf("[%s] %s\n", url.Host, err.Error())
		return false
	}

	defer conn.Close()

	if err != nil {
		log.Println("Site unreachable , error:", err)
		return false
	}
	return true
}
