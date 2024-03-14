package deststream

import (
	"buffer-handler/logging"
	"buffer-handler/models"
)

type KafkaStreamConfig struct {
	BootstrapServers []string
	Topic            string
}

func (cfg KafkaStreamConfig) Notify(models.EntryData) error {
	logging.Logger.Debug("Kafka Streaming")

	return nil
}
