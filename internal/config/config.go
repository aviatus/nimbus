package config

import (
	"log"
	"os"

	lb "github.com/aviatus/nimbus/internal/loadbalancer"
	"gopkg.in/yaml.v2"
)

type Config struct {
	LoadBalancer lb.LoadBalancerConfig `yaml:"loadBalancer"`
}

func ImportConfig() *lb.LoadBalancer {
	// Read the YAML configuration file
	configData, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}

	// Parse the YAML data into a Config struct
	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf("Error parsing configuration data: %v", err)
	}

	loadBalancer := lb.NewLoadBalancer(config.LoadBalancer)
	return loadBalancer
}
