package loadbalancer

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/aviatus/nimbus/internal/server"
	"github.com/aviatus/nimbus/internal/service"
)

type LoadBalancerConfig struct {
	Port           string                  `yaml:"port"`
	ManagementPort string                  `yaml:"managementPort"`
	Timeout        time.Duration           `yaml:"timeout"`
	Services       []service.ServiceConfig `yaml:"services"`
	TLSConfigs     TLSConfigs              `yaml:"tlsConfig"`
}

type LoadBalancer struct {
	Services       []*service.Service
	Proxy          *httputil.ReverseProxy
	Port           string
	ManagementPort string
	Timeout        time.Duration
	Server         *http.Server
	TLSConfigs     TLSConfigs
}

type TLSConfigs struct {
	Enabled  bool   `yaml:"enabled"`
	Port     string `yaml:"port"`
	CertPath string `yaml:"certPath"`
	KeyPath  string `yaml:"keyPath"`
}

func NewLoadBalancer(config LoadBalancerConfig) *LoadBalancer {
	services := []*service.Service{}
	for _, serviceConfig := range config.Services {
		service := service.NewService(serviceConfig)
		services = append(services, service)
	}

	lb := &LoadBalancer{
		Services:   services,
		Proxy:      &httputil.ReverseProxy{},
		Port:       config.Port,
		Timeout:    config.Timeout,
		TLSConfigs: config.TLSConfigs,
	}
	lb.Proxy.Director = lb.director

	return lb
}

func StartLoadBalancer(lb *LoadBalancer) {
	port := lb.Port
	scheme := "http"
	if lb.TLSConfigs.Enabled {
		port = lb.TLSConfigs.Port
		scheme = "https"
	}
	lb.Server = &http.Server{
		Addr: ":" + port,
		Handler: http.TimeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Scheme = scheme
			r.URL.Host = r.Host
			lb.Proxy.ServeHTTP(w, r)
		}), lb.Timeout*time.Second, "Timeout"),
		WriteTimeout: lb.Timeout * time.Second,
	}
	
	fmt.Printf("Load balancer listening on port %s\n", port)
	switch lb.TLSConfigs.Enabled {
	case true:
		go func() {
			if err := lb.Server.ListenAndServeTLS(lb.TLSConfigs.CertPath, lb.TLSConfigs.KeyPath); err != nil {
				panic(err)
			}
		}()

		if lb.Port != "" {
			httpServer := &http.Server{
				Addr:         ":" + lb.Port,
				Handler:      http.HandlerFunc(redirectToHTTPS),
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
			}
			go func() {
				fmt.Printf("HTTP server listening on port %s, redirecting to HTTPS\n", lb.Port)
				if err := httpServer.ListenAndServe(); err != nil {
					panic(err)
				}
			}()
		}
	default:
		go func() {
			if err := lb.Server.ListenAndServe(); err != nil {
				panic(err)
			}
		}()
	}
}

func (lb *LoadBalancer) director(req *http.Request) {
	servers := []*server.Server{}
	for _, service := range lb.Services {
		if service.Host == req.Host {
			servers = append(servers, service.Servers...)
			break
		} 
	}

	if len(servers) == 0 {
		fmt.Printf("No servers found for host %s\n", req.Host)
		return
	}

	maxAttempts := len(servers)
	for i := 0; i < maxAttempts; i++ {
		serverIndex := rand.Intn(maxAttempts)
		server := servers[serverIndex]

		if server.IsAlive() {
			select {
			case server.ConnectionCh <- struct{}{}:
				req.URL.Host = server.URL
				server.HandleRequest(req)
				return
			default:
				server.QueueCh <- req
			}
			break
		}
	}
}

func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	httpsURL := "https://" + r.Host + r.URL.String()
	http.Redirect(w, r, httpsURL, http.StatusPermanentRedirect)
}

func (lb *LoadBalancer) UpdateServices(services []service.ServiceConfig) {
	// newServers := make(map[string]*service.Service)
	// for _, serviceConfig := range services {
	// 	newServer := server.NewServer()
	// 	newServers[newServer.Id] = newServer
	// }

	// for i := len(lb.Services) - 1; i >= 0; i-- {
	// 	if _, exists := newServices[lb.Services[i].ID]; !exists {
	// 		lb.Services = append(lb.Services[:i], lb.Services[i+1:]...)
	// 	}
	// }

	// for _, newService := range newServices {
	// 	exists := false
	// 	for _, existingService := range lb.Services {
	// 		if existingService.ID == newService.ID {
	// 			exists = true
	// 			break
	// 		}
	// 	}
	// 	if !exists {
	// 		lb.Services = append(lb.Services, newService)
	// 	}
	// }
}
