package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aviatus/nimbus/internal/config"
	"github.com/aviatus/nimbus/internal/loadbalancer"
	mng "github.com/aviatus/nimbus/internal/management"
	"github.com/aviatus/nimbus/internal/utils"
)

func main() {
	debugFlag := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	utils.SetDebugEnabled(*debugFlag)

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}
	config.SetConfig(cfg)

	lb := loadbalancer.NewLoadBalancer(cfg.LoadBalancer)
	mng.StartManagementServer(lb)
	lb.StartLoadBalancer()

	// Set up channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-sigChan

	fmt.Println("Interrupt signal received, shutting down...")

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := lb.Shutdown(ctx); err != nil {
		log.Fatalf("Nimbus shutdown failed: %v", err)
	}

	log.Println("Nimbus shut down gracefully")
}
