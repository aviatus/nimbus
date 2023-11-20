package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()
	// Create a mock HTTP server for testing
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a healthy response (HTTP status code 200)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	// Create a server instance with the mock URL
	server := NewServer(mockServer.URL, mockServer.URL+"/health", 1)

	// Start the health check in a goroutine
	go server.HealthCheck()

	// Sleep for a while to allow the HealthCheck goroutine to run
	time.Sleep(6 * time.Second)

	// Check if the server's "Alive" field is set to true (indicating it's healthy)
	if !server.IsAlive() {
		t.Error("Expected server to be alive, but it's not.")
	}

	// You can add more test cases to simulate different scenarios, such as an unhealthy server response.
}

// func TestServer_HandleQueue(t *testing.T) {
// 	t.Parallel()
// 	// Create a server instance with a concurrency limit of 2 and a queue buffer of 1
// 	server := NewServer("http://example.com", "http://example.com/health", 2)
// 	server.QueueCh = make(chan *http.Request, 1)

// 	// Simulate sending two requests to the server
// 	request1 := &http.Request{}
// 	request2 := &http.Request{}
// 	server.QueueCh <- request1
// 	server.QueueCh <- request2

// 	// Start the HandleQueue goroutine in the background
// 	go server.HandleQueue()

// 	// Allow some time for the HandleQueue goroutine to run
// 	time.Sleep(1 * time.Second)

// 	// Check if the first request was handled (should be removed from ConnectionCh)
// 	select {
// 	case <-server.ConnectionCh:
// 	default:
// 		t.Error("Expected the first request to be handled, but it was not.")
// 	}

// 	// Check if the second request was kept in the queue
// 	select {
// 	case <-server.QueueCh:
// 		t.Error("Expected the second request to be in the queue, but it was not.")
// 	default:
// 		// The request is still in the queue, which is expected
// 	}
// }
