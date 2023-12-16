package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aviatus/nimbus/internal/utils"
)

func NewServer(cfg ServerConfig) *Server {
	server := &Server{
		URL:              cfg.URL,
		HealthURL:        cfg.HealthURL,
		HealthStatusChan: make(chan bool),
		alive:            false,
	}
	hash, err := server.CreateHash(cfg)
	if err != nil {
		fmt.Println("Error creating hash for server")
	}
	server.hash = hash
	return server
}

func (s *Server) HealthCheck() {
	fmt.Println("Health checking server: ", s.URL)

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	for {
		time.Sleep(5 * time.Second)
		resp, err := client.Get(s.HealthURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			if s.IsAlive() {
				s.HealthStatusChan <- false
				s.SetAlive(false)
				utils.Debugf("Server %s is dead\n", s.URL)
			}
		} else if !s.IsAlive() {
			s.SetAlive(true)
			s.HealthStatusChan <- true
			resp.Body.Close()
			utils.Debugf("Server %s is alive\n", s.URL)
		}
	}
}

func (s *Server) CreateHash(cfg ServerConfig) (string, error) {
	data := fmt.Sprintf("%s|%s", s.URL, s.HealthURL)
	return utils.HashObject(data)
}

func (s *Server) CompareHash(hash string) bool {
	return s.hash == hash
}

func (s *Server) IsAlive() bool {
	return s.alive
}

func (s *Server) SetAlive(alive bool) {
	s.alive = alive
}

func (s *Server) Shutdown() {
	s.SetAlive(false)
}
