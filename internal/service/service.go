package service

import (
	"github.com/aviatus/nimbus/internal/server"
)

type ServiceConfig struct {
	Name             string                `yaml:"name"`
	Host             string                `yaml:"host"`
	ConcurrencyLimit int                   `yaml:"concurrencyLimit"`
	Servers          []server.ServerConfig `yaml:"servers"`
}

type Service struct {
	Name             string
	Host             string
	ConcurrencyLimit int
	Servers          []*server.Server
}

func NewService(config ServiceConfig) *Service {
	var servers []*server.Server
	for _, serverConfig := range config.Servers {
		server := server.NewServer(serverConfig.URL, serverConfig.HealthURL, config.ConcurrencyLimit)
		servers = append(servers, server)
		go server.HealthCheck()
		go server.HandleQueue()
	}

	service := &Service{
		Name:             config.Name,
		Host:             config.Host,
		ConcurrencyLimit: config.ConcurrencyLimit,
		Servers:          servers,
	}

	return service
}
