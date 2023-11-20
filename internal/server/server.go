package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/aviatus/nimbus/internal/utils"
)

type ServerConfig struct {
	URL       string `yaml:"url"`
	HealthURL string `yaml:"healthURL"`
}

type Server struct {
	Id		   	   string
	URL            string
	Alive          bool
	Mutex          sync.Mutex
	HealthURL      string
	ConcurrencyLim int
	ConnectionCh   chan struct{}
	QueueCh        chan *http.Request
}

func NewServer(url, healthURL string, concurrencyLim int) *Server {
	server := &Server{
		URL:            url,
		HealthURL:      healthURL,
		Alive:          false,
		ConcurrencyLim: concurrencyLim,
		ConnectionCh:   make(chan struct{}, concurrencyLim),
		QueueCh:        make(chan *http.Request, 1000),
	}
	server.Id, _ = utils.HashObject(&server)
	return server
}

func (s *Server) HealthCheck() {
	for {
		time.Sleep(5 * time.Second)
		resp, err := http.Get(s.HealthURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			s.SetAlive(false)
			fmt.Printf("Server %s is dead\n", s.URL)
		} else {
			s.SetAlive(true)
			resp.Body.Close()
			fmt.Printf("Server %s is alive\n", s.URL)
		}
	}
}

func (s *Server) HandleQueue() {
	for req := range s.QueueCh {
		select {
		case s.ConnectionCh <- struct{}{}:
			// Server has room for more connections
			s.HandleRequest(req)
		default:
			// Server is still at concurrency limit, keep it in the queue
			s.QueueCh <- req
		}
	}
}

func (s *Server) IsAlive() bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.Alive
}

func (s *Server) SetAlive(alive bool) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Alive = alive
}

func (s *Server) HandleRequest(req *http.Request) {
	defer func() {
		<-s.ConnectionCh
	}()
}

func (s *Server) Shutdown() {
    // Signal to stop accepting new requests
    close(s.QueueCh)

    // Wait for all connections to be processed
    for range s.ConnectionCh { 
		// This loop will continue until the ConnectionCh is closed and emptied
	}
}
