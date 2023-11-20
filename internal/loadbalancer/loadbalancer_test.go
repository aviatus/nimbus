package loadbalancer

import (
	"net/http"
	"testing"
	"time"

	"github.com/aviatus/nimbus/internal/server"
	"github.com/aviatus/nimbus/internal/service"
)

func TestNewLoadBalancer(t *testing.T) {
	t.Parallel()
	// Create a mock configuration
	mockConfig := LoadBalancerConfig{
		Port:           "8080",
		ManagementPort: "8081",
		Timeout:        10 * time.Second,
		Services: []service.ServiceConfig{
			{
				// Add service configurations here
			},
			// Add more service configurations as needed
		},
	}

	// Create a new LoadBalancer using the mock configuration
	lb := NewLoadBalancer(mockConfig)

	// Add assertions to check if the LoadBalancer is initialized correctly
	if lb.Port != "8080" {
		t.Errorf("Expected Port to be '8080', but got '%s'", lb.Port)
	}

	if lb.Timeout != 10*time.Second {
		t.Errorf("Expected Timeout to be 10 seconds, but got '%s'", lb.Timeout)
	}

	// Add more assertions to check other fields and configurations

	// For example, you can check if lb.Services is correctly initialized
	if len(lb.Services) != len(mockConfig.Services) {
		t.Errorf("Expected %d services, but got %d", len(mockConfig.Services), len(lb.Services))
	}
}

func TestLoadBalancer_StartLoadBalancer(t *testing.T) {
	t.Parallel()
	// Create a mock LoadBalancer instance with configuration
	mockConfig := LoadBalancerConfig{
		Port:           "8080",
		ManagementPort: "8081",
		Timeout:        5,   // 5 seconds
		Services:       nil, // You can define services if needed
	}
	mockLB := NewLoadBalancer(mockConfig)

	// Start the load balancer in the background
	StartLoadBalancer(mockLB)
	go func() {
		// Simulate an HTTP request to the load balancer
		requestURL := "http://localhost:8080/somepath"
		resp, err := http.Get(requestURL)
		if err != nil {
			t.Errorf("Failed to make a request: %v", err)
		}

		// Verify that the response status code is as expected
		if resp.StatusCode != http.StatusGatewayTimeout {
			t.Errorf("Expected status code %d, but got %d", http.StatusGatewayTimeout, resp.StatusCode)
		}
	}()
}

func TestLoadBalancer_Director(t *testing.T) {
	// Create a LoadBalancer instance with a mock set of services
	lb := &LoadBalancer{
		Services: []*service.Service{
			{
				Host: "example.com",
				Servers: []*server.Server{
					{
						URL:          "http://server1",
						Alive:        true,
						ConnectionCh: make(chan struct{}),
					},
					{
						URL:   "http://server2",
						Alive: false,
					},
				},
			},
		},
	}

	// Create a mock HTTP request for testing
	req := &http.Request{
		Host: "example.com",
	}

	go func() {
		// Call the director method to make routing decisions
		lb.director(req)

		// For example, you can assert that req.URL.Host has been set to a valid server URL
		if req.URL.Host != "http://server1" {
			t.Errorf("Expected req.URL.Host to be 'http://server1', but got '%s'", req.URL.Host)
		}
	}()
}
