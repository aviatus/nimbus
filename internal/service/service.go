package service

import (
	"github.com/aviatus/nimbus/internal/server"
)

func NewService(config ServiceConfig) *Service {
	var servers []*server.Server
	for _, serverConfig := range config.Servers {
		server := server.NewServer(*serverConfig)
		servers = append(servers, server)
		go server.HealthCheck()
	}

	healthCheckInterval := config.HealthCheck.HealthCheckInterval
	if healthCheckInterval == 0 {
		healthCheckInterval = 10 // default value in seconds
	}

	healthCheckTimeout := config.HealthCheck.HealthCheckTimeout
	if healthCheckTimeout == 0 {
		healthCheckTimeout = 5 // default value in seconds
	}

	service := &Service{
		Name:    config.Name,
		Host:    config.Host,
		servers: servers,
		healthCheck: HealthCheckConfig{
			HealthCheckInterval: healthCheckInterval,
			HealthCheckTimeout:  healthCheckTimeout,
		},
	}

	go service.monitorHealth()
	return service
}

func (s *Service) GetServers() []*server.Server {
	return s.servers
}

func (s *Service) GetAvailableServers() []*server.Server {
	var availableServers []*server.Server
	for _, sv := range s.servers {
		if sv.IsAlive() {
			availableServers = append(availableServers, sv)
		}
	}
	return availableServers
}

func (s *Service) IsAlive() bool {
	return s.availableServerCount > 0
}

func (s *Service) monitorHealth() {
	for _, sv := range s.servers {
		go func(ser *server.Server) {
			for {
				select {
				case status := <-ser.HealthStatusChan:
					if status {
						s.availableServerCount++
					} else {
						s.availableServerCount--
					}
				}
			}
		}(sv)
	}
}

func (s *Service) UpdateName(name string) {
	s.Name = name
}

func (s *Service) UpdateHost(host string) {
	s.Host = host
}

func (s *Service) UpdateServers(servers []server.ServerConfig) {

}
