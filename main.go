package main

import (
	"github.com/aviatus/nimbus/internal/config"
	"github.com/aviatus/nimbus/internal/loadbalancer"
	"github.com/aviatus/nimbus/internal/management"
)

func main() {
	config := config.ImportConfig()
	management.StartManagementServer(config)
	loadbalancer.StartLoadBalancer(config)

	select {}
}
