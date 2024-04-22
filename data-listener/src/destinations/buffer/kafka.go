package buffer

import (
	"datalistener/src/logging"
	"datalistener/src/models"
)

type KafkaBufferConfig struct {
	BootstrapServers []string
	Topic            string
}

func (cfg KafkaBufferConfig) Notify(entries []models.EntryData, convertToJSONL bool) error {
	logging.Logger.Debug("Kafka Buffering")

	return nil
}
