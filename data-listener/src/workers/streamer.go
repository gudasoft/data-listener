package workers

import (
	"datalistener/src/destinations"
	"datalistener/src/logging"
	"datalistener/src/models"
	"reflect"
	"sync"
)

func DestinationsNotifier(streamerConfigs *[]destinations.StreamConfig, streamerChannel chan models.EntryData,
	shutdownChannel chan struct{}, readyChannel chan bool) {
	wg := sync.WaitGroup{}

selectloop:
	for {
		select {
		case <-shutdownChannel:

			wg.Wait()
			readyChannel <- true

			break selectloop
		case data := <-streamerChannel:
			for _, destination := range *streamerConfigs {
				wg.Add(1)
				go func(d destinations.StreamConfig) {
					defer wg.Done()
					if err := d.Notify(data); err != nil {
						destinationType := reflect.TypeOf(d).String()
						logging.Logger.Sugar().Errorf("Error notifying streamer %s, %s, %v", err, destinationType, "// TODO tmp")
					}
				}(destination)
			}
			logging.Logger.Sugar().Debugf("Streamer processed: %s", data)
		}
	}
}
