package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/brianlusina/robin-lb/app/api/lb"
	"github.com/brianlusina/robin-lb/app/internal/services/proxysrv"
	"github.com/brianlusina/robin-lb/app/internal/services/robinsrv"
)

func main() {
	var serverList string
	var port int
	flag.StringVar(&serverList, "backends", "", "Load balanced backends, use commas to separate")
	flag.IntVar(&port, "port", 3030, "Port to serve")
	flag.Parse()

	if len(serverList) == 0 {
		log.Fatal("Please provide one or more backends to load balance")
	}

	servers := strings.Split(serverList, ",")
	serverPool := robinsrv.New()
	lb := lb.NewHandler(serverPool)

	for _, server := range servers {
		serverURL, err := url.Parse(server)
		if err != nil {
			log.Fatalf("Failed to parse server URL: %s", err)
		}

		proxySvc := proxysrv.New(serverURL, lb.LoadBalance, serverPool)
		proxySvc.AddErrorHandler()

		serverPool.AddBackend(serverURL, true, proxySvc)

		log.Printf("Configured server: %s\n", serverURL)
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(lb.LoadBalance),
	}

	// performs a healthcheck on the server pool to check on status of backends every 2 mins
	go func() {
		t := time.NewTicker(time.Minute * 2)
		for range t.C {
			log.Println("Starting health check...")
			serverPool.HealthCheck()
			log.Println("Health check completed")
		}
	}()

	log.Printf("Robin Load Balancer started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
