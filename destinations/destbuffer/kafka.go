package destbuffer

import (
	"buffer-handler/logging"
	"buffer-handler/models"
)

type KafkaBufferConfig struct {
	BootstrapServers []string
	Topic            string
}

func (cfg KafkaBufferConfig) Notify(entries []models.EntryData, convertToJSONL bool) error {
	logging.Logger.Debug("Kafka Buffering")

	return nil
}
