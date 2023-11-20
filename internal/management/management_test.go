package management

import (
	"net/http"
	"testing"

	lb "github.com/aviatus/nimbus/internal/loadbalancer"
)

func TestStartManagementServer(t *testing.T) {
	// Create a mock LoadBalancer instance
	mockLB := &lb.LoadBalancer{
		ManagementPort: "8080",
	}

	// Start the management server in the background
	go StartManagementServer(mockLB)

	// Simulate an HTTP request to the /reloadConfig endpoint
	reloadConfigURL := "http://localhost:8080/reloadConfig"
	resp := sendRequest(reloadConfigURL)

	// Verify the response from the /reloadConfig endpoint
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200 for /reloadConfig, but got %d", resp.StatusCode)
	}

	// Simulate an HTTP request to a non-existent endpoint
	nonExistentURL := "http://localhost:8080/nonexistent"
	resp = sendRequest(nonExistentURL)

	// Verify the response from the non-existent endpoint
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404 for /nonexistent, but got %d", resp.StatusCode)
	}
}

func sendRequest(url string) *http.Response {
	// Send an HTTP GET request to the specified URL
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	return resp
}
