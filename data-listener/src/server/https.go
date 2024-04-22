package server

import "fmt"

type HttpsServerConfig struct {
	Protocol       string
	Address        string
	Port           int
	ServerCertFile string
	ServerKeyFile  string
	UseMTLS        bool
	CACertFile     string
}

func (cfg HttpsServerConfig) Start() string {
	return "KEEPIT"
}

func (cfg HttpsServerConfig) String() string {
	if cfg.UseMTLS {
		return fmt.Sprintf("MTLS  %s://%s:%d", cfg.Protocol, cfg.Address, cfg.Port)
	}
	return fmt.Sprintf("%s://%s:%d", cfg.Protocol, cfg.Address, cfg.Port)
}
