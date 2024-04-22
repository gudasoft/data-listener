package stream

import (
	"datalistener/src/logging"
	"datalistener/src/models"
)

type KafkaStreamConfig struct {
	BootstrapServers []string
	Topic            string
}

func (cfg KafkaStreamConfig) Notify(models.EntryData) error {
	logging.Logger.Debug("Kafka Streaming")

	return nil
}
