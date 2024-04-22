package stream

import (
	"datalistener/src/logging"
	"datalistener/src/models"
)

type SmtpsStreamConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	UseTLS   bool
}

func (cfg *SmtpsStreamConfig) Notify(models.EntryData) error {
	logging.Logger.Debug("Smtps Streaming")

	return nil
}
