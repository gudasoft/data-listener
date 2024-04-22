package server

type UnixServerConfig struct {
	Protocol string
	Address  string
}

func (cfg UnixServerConfig) Start() string {
	return "KEEPIT"
}
