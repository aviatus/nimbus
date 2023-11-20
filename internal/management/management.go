package management

import (
	"fmt"
	"net/http"

	lb "github.com/aviatus/nimbus/internal/loadbalancer"
)

const defaultPort = "10024"

func StartManagementServer(lb *lb.LoadBalancer) {
	// First HTTP server on port 8080
	http.HandleFunc("/reloadConfig", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", r.URL.Path)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	port := defaultPort
	if lb.ManagementPort != "" {
		port = lb.ManagementPort
	}

	go func() {
		fmt.Println("Management server listening on port: ", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			panic(err)
		}
	}()
}
