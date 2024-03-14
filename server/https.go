package server

import "fmt"

type HttpsServerConfig struct {
	Protocol string
	Address  string
	Port     int
	TlsCert  string
	TlsKey   string
}

func (cfg HttpsServerConfig) Start() string {
	return "KEEPIT"
}

func (cfg HttpsServerConfig) String() string {
	return fmt.Sprintf("%s://%s:%d", cfg.Protocol, cfg.Address, cfg.Port)
}
