package server

type ServerConfig struct {
	URL       string `yaml:"url"`
	HealthURL string `yaml:"healthURL"`
}

type Server struct {
	URL              string
	HealthURL        string
	HealthStatusChan chan bool
	hash             string
	alive            bool
}
