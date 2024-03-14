package server

import "fmt"

type PrometheusServerConfig struct {
	Address string
	Port    int
	Path    string
}

func (cfg PrometheusServerConfig) Start() string {
	return "KEEPIT"
}

func (cfg PrometheusServerConfig) String() string {
	return fmt.Sprintf("http://unix:%s:%d%s", cfg.Address, cfg.Port, cfg.Path)
}
