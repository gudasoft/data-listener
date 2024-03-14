package workers

import (
	"buffer-handler/destinations"
	"buffer-handler/logging"
	"buffer-handler/models"
)

func DestinationsNotifier(streamerConfigs *[]destinations.StreamConfig, streamerChannel chan models.EntryData) error {

	for data := range streamerChannel {
		for _, destionation := range *streamerConfigs {
			destionation.Notify(data)
		}
		logging.Logger.Sugar().Info("Streamer processed:", data)
	}

	return nil
}
