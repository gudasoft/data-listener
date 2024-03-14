package server

import "fmt"

type HttpServerConfig struct {
	Protocol string
	Unix     string
	Address  string
	Port     int
	CertFile string
	KeyFile  string
}

func (cfg HttpServerConfig) Start() string {
	return "KEEPIT"
}

func (cfg HttpServerConfig) String() string {
	if cfg.Protocol == "http" {
		return fmt.Sprintf("%s://%s:%d", cfg.Protocol, cfg.Address, cfg.Port)
	}
	return fmt.Sprintf("HttpServerConfig{Protocol: %s, Unix: %s, Address: %s, Port: %d}",
		cfg.Protocol, cfg.Unix, cfg.Address, cfg.Port)
}
