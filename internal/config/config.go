package config

import (
	"log"
	"os"

	lb "github.com/aviatus/nimbus/internal/loadbalancer"
	"github.com/aviatus/nimbus/internal/server"
	"github.com/aviatus/nimbus/internal/service"
	"gopkg.in/yaml.v2"
)

type Config struct {
	LoadBalancer lb.LoadBalancerConfig `yaml:"loadBalancer"`
}

var config *Config

func SetConfig(c *Config) {
	config = c
}

func GetConfig() *Config {
	return config
}

func LoadConfig(filename string) (*Config, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func RefreshConfig() {
	newConfig, err := LoadConfig("config.yaml")
	if err != nil {
		log.Printf("Error reloading config: %v", err)
	} else {
		UpdateLoadBalancerConfig(&config.LoadBalancer, newConfig.LoadBalancer)
		SetConfig(newConfig)
	}
}

// UpdateServerConfig updates the server configuration if changes are detected
func UpdateServerConfig(current *server.ServerConfig, new server.ServerConfig) {
	// Update fields only if they have changed
	if current.URL != new.URL {
		current.URL = new.URL
	}
	if current.HealthURL != new.HealthURL {
		current.HealthURL = new.HealthURL
	}
	// Add checks for other fields if necessary
}

// UpdateServiceConfig updates the service configuration and its servers
func UpdateServiceConfig(current *service.ServiceConfig, new service.ServiceConfig) {
	if current.Host != new.Host {
		current.Host = new.Host
	}
	// Update servers
	for _, newServer := range new.Servers {
		// Find the matching server in the current config
		var found bool
		for i := range current.Servers {
			if current.Servers[i].URL == newServer.URL { // Assuming URL is a unique identifier
				UpdateServerConfig(current.Servers[i], *newServer)
				found = true
				break
			}
		}
		// If the server is new, add it to the current config
		if !found {
			current.Servers = append(current.Servers, newServer)
		}
	}
	// Optionally handle removal of servers no longer in the new config
}

// UpdateLoadBalancerConfig updates the load balancer configuration and its services
func UpdateLoadBalancerConfig(current *lb.LoadBalancerConfig, new lb.LoadBalancerConfig) {
	// if current.Port != new.Port {
	// 	current.Port = new.Port
	// }

	// Update services
	for _, newService := range new.Services {
		var found bool
		for i := range current.Services {
			if current.Services[i].Name == newService.Name { // Assuming Name is a unique identifier
				UpdateServiceConfig(&current.Services[i], newService)
				found = true
				break
			}
		}
		if !found {
			current.Services = append(current.Services, newService)
		}
	}
	// Optionally handle removal of services no longer in the new config
}
