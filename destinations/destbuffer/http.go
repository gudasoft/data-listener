package destbuffer

import (
	"buffer-handler/logging"
	"buffer-handler/models"
	"fmt"
)

type HttpBufferConfig struct {
	Protocol string
	Address  string
	Port     int
}

func (cfg HttpBufferConfig) Notify(entries []models.EntryData, convertToJSONL bool) error {
	logging.Logger.Debug("Http Buffering")
	return nil
}

func (cfg HttpBufferConfig) String() string {
	return fmt.Sprintf("%s://%s:%d", cfg.Protocol, cfg.Address, cfg.Port)
}
