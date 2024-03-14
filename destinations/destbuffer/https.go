package destbuffer

import (
	"buffer-handler/logging"
	"buffer-handler/models"
	"fmt"
)

type HttpsBufferConfig struct {
	Protocol string
	Address  string
	Port     int
	TlsCert  string
	TlsKey   string
}

func (cfg HttpsBufferConfig) Notify(entries []models.EntryData, convertToJSONL bool) error {
	logging.Logger.Debug("Https Buffering")

	return nil
}

func (cfg HttpsBufferConfig) String() string {
	return fmt.Sprintf("%s://%s:%d", cfg.Protocol, cfg.Address, cfg.Port)
}
