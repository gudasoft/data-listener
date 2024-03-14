package deststream

import (
	"buffer-handler/logging"
	"buffer-handler/models"
	"fmt"
)

type HttpStreamConfig struct {
	Protocol string
	Address  string
	Port     int
}

func (cfg HttpStreamConfig) Notify(data models.EntryData) error {
	logging.Logger.Debug("Http Streaming")
	return nil
}

func (cfg HttpStreamConfig) String() string {
	return fmt.Sprintf("%s://%s:%d", cfg.Protocol, cfg.Address, cfg.Port)
}
