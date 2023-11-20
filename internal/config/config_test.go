package config

import (
	"os"
	"testing"
)

func TestImportConfig(t *testing.T) {
	// Create a temporary test configuration file
	const testConfig = `
loadBalancer:
  port: "8080"
  managementPort: "8081"
  timeout: 5
  services:
    - name: "Service1"
      host: "example.com"
      concurrencyLimit: 10
      servers:
        - url: "http://server1.com"
          healthURL: "http://server1.com/health"
    - name: "Service2"
      host: "example.org"
      concurrencyLimit: 5
      servers:
        - url: "http://server2.com"
          healthURL: "http://server2.com/health"
`
	err := os.WriteFile("config.yaml", []byte(testConfig), 0666)
	if err != nil {
		t.Fatalf("Failed to create test configuration file: %v", err)
	}
	defer os.Remove("config.yaml")

	// Test ImportConfig
	loadBalancer := ImportConfig()

	// Verify that the LoadBalancer configuration is parsed correctly
	if loadBalancer.Port != "8080" {
		t.Errorf("Expected Port to be '8080', but got '%s'", loadBalancer.Port)
	}

	// Add more assertions to verify other configuration properties
}
