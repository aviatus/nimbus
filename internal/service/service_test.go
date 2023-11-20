package service

import (
	"testing"

	"github.com/aviatus/nimbus/internal/server"
)

func TestNewService(t *testing.T) {
	// Create a mock ServiceConfig for testing
	mockConfig := ServiceConfig{
		Name:             "TestService",
		Host:             "example.com",
		ConcurrencyLimit: 10,
		Servers: []server.ServerConfig{
			{
				URL:       "http://server1.com",
				HealthURL: "http://server1.com/health",
			},
			{
				URL:       "http://server2.com",
				HealthURL: "http://server2.com/health",
			},
		},
	}

	// Create a new Service using the NewService function
	service := NewService(mockConfig)

	// Verify that the service is created correctly
	if service.Name != mockConfig.Name {
		t.Errorf("Expected service name to be %s, but got %s", mockConfig.Name, service.Name)
	}

	if service.Host != mockConfig.Host {
		t.Errorf("Expected service host to be %s, but got %s", mockConfig.Host, service.Host)
	}

	if service.ConcurrencyLimit != mockConfig.ConcurrencyLimit {
		t.Errorf("Expected concurrency limit to be %d, but got %d", mockConfig.ConcurrencyLimit, service.ConcurrencyLimit)
	}

	// Add more assertions to check the correctness of the created service object

	// You can also check if the servers slice is created correctly based on the ServerConfig objects

	if len(service.Servers) != len(mockConfig.Servers) {
		t.Errorf("Expected %d servers, but got %d", len(mockConfig.Servers), len(service.Servers))
	}

	// Perform additional assertions on the server instances, if needed

	// Clean up or defer any necessary teardown
}
