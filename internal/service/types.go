package service

import "github.com/aviatus/nimbus/internal/server"

type ServiceConfig struct {
	Name        string                 `yaml:"name"`
	Host        string                 `yaml:"host"`
	Servers     []*server.ServerConfig `yaml:"servers"`
	HealthCheck HealthCheckConfig      `yaml:"healthCheck"`
}

type Service struct {
	Name                 string
	Host                 string
	servers              []*server.Server
	healthCheck          HealthCheckConfig
	availableServerCount int
	hash                 string
}

type HealthCheckConfig struct {
	HealthCheckInterval int `yaml:"healthCheckInterval"`
	HealthCheckTimeout  int `yaml:"healthCheckTimeout"`
}
