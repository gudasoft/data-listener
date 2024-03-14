package deststream

import (
	"buffer-handler/logging"
	"buffer-handler/models"
	"fmt"
)

type HttpsStreamConfig struct {
	Protocol string
	Address  string
	Port     int
	TlsCert  string
	TlsKey   string
}

func (cfg HttpsStreamConfig) Notify(entry models.EntryData) error {
	logging.Logger.Debug("Https Streaming")

	return nil
}

func (cfg HttpsStreamConfig) String() string {
	return fmt.Sprintf("%s://%s:%d", cfg.Protocol, cfg.Address, cfg.Port)
}
