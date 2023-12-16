package management

import (
	"fmt"
	"net/http"

	"github.com/aviatus/nimbus/internal/config"
	"github.com/aviatus/nimbus/internal/loadbalancer"
	"github.com/aviatus/nimbus/internal/utils"
)

const defaultManagementPort = "10024"

func StartManagementServer(lb *loadbalancer.LoadBalancer) {
	http.HandleFunc("/reloadConfig", func(w http.ResponseWriter, r *http.Request) {
		config.RefreshConfig()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	port := defaultManagementPort
	if lb.ManagementPort != "" {
		port = lb.ManagementPort
	}

	go func() {
		fmt.Println("Management server listening on port: ", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			utils.Debugln(err)
		}
	}()
}
