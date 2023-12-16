package loadbalancer

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/aviatus/nimbus/internal/service"
)

type LoadBalancerConfig struct {
	Port                 string                  `yaml:"port"`
	ManagementPort       string                  `yaml:"managementPort"`
	Timeout              time.Duration           `yaml:"timeout"`
	Services             []service.ServiceConfig `yaml:"services"`
	TLSConfigs           TLSConfigs              `yaml:"tlsConfig"`
	ConnectionPoolConfig ConnectionPoolConfig    `yaml:"connectionPoolConfig"`
}

type ConnectionPoolConfig struct {
	Enabled             bool          `yaml:"enabled"`
	MaxIdleConns        int           `yaml:"maxIdleConns"`
	MaxIdleConnsPerHost int           `yaml:"maxIdleConnsPerHost"`
	IdleConnTimeout     time.Duration `yaml:"idleConnTimeout"`
}

type LoadBalancer struct {
	Services       []*service.Service
	Proxy          *httputil.ReverseProxy
	Port           string
	ManagementPort string
	Timeout        time.Duration
	Server         *http.Server
	TLSConfigs     TLSConfigs
	Transport      *http.Transport
	hash           string
}

type TLSConfigs struct {
	Enabled  bool   `yaml:"enabled"`
	Port     string `yaml:"port"`
	CertPath string `yaml:"certPath"`
	KeyPath  string `yaml:"keyPath"`
}
