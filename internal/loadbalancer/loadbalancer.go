package loadbalancer

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/aviatus/nimbus/internal/service"
	"github.com/aviatus/nimbus/internal/utils"
)

func NewLoadBalancer(config LoadBalancerConfig) *LoadBalancer {
	services := []*service.Service{}
	for _, serviceConfig := range config.Services {
		service := service.NewService(serviceConfig)
		services = append(services, service)
	}

	lb := &LoadBalancer{
		Services:       services,
		Proxy:          &httputil.ReverseProxy{},
		Port:           config.Port,
		Timeout:        config.Timeout,
		TLSConfigs:     config.TLSConfigs,
		ManagementPort: config.ManagementPort,
	}

	lb.Proxy.Director = lb.director
	if config.ConnectionPoolConfig.Enabled {
		lb.Transport = &http.Transport{
			MaxIdleConns:        config.ConnectionPoolConfig.MaxIdleConns,
			MaxIdleConnsPerHost: config.ConnectionPoolConfig.MaxIdleConnsPerHost,
			IdleConnTimeout:     config.ConnectionPoolConfig.IdleConnTimeout,
		}
		lb.Proxy.Transport = lb.Transport
	}

	return lb
}

func (lb *LoadBalancer) StartLoadBalancer() {
	port := lb.Port
	scheme := "http"
	if lb.TLSConfigs.Enabled {
		port = lb.TLSConfigs.Port
		// scheme = "https"
	}

	lb.Server = &http.Server{
		Addr: ":" + port,
		Handler: http.TimeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			service := lb.getService(r.Host)
			if service == nil || !service.IsAlive() {
				http.Error(w, "Domain is not alive", http.StatusServiceUnavailable)
				return
			}

			r.URL.Scheme = scheme
			r.URL.Host = r.Host
			lb.Proxy.ServeHTTP(w, r)
		}), lb.Timeout*time.Second, "Timeout"),
		WriteTimeout: lb.Timeout * time.Second,
		IdleTimeout:  lb.Timeout * time.Second,
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
					utils.Debugln(err)
				}
			}()
		}
	default:
		go func() {
			if err := lb.Server.ListenAndServe(); err != nil {
				utils.Debugln(err)
			}
		}()
	}
}

func (lb *LoadBalancer) director(req *http.Request) {
	servers := lb.getService(req.Host).GetAvailableServers()
	serverIndex := rand.Intn(len(servers))

	req.URL.Host = servers[serverIndex].URL
}

func (lb *LoadBalancer) getService(name string) *service.Service {
	for _, service := range lb.Services {
		if service.Host == name {
			return service
		}
	}

	return nil
}

func (lb *LoadBalancer) Shutdown(ctx context.Context) error {
	return lb.Server.Shutdown(ctx)
}

func (lb *LoadBalancer) CloseIdleConnections() {
	lb.Transport.CloseIdleConnections()
}

func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	httpsURL := "https://" + r.Host + r.URL.String()
	http.Redirect(w, r, httpsURL, http.StatusPermanentRedirect)
}
