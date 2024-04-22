package workers

import (
	"datalistener/src/destinations"
	"datalistener/src/logging"
	"datalistener/src/models"
	"reflect"
	"sync"
)

func DestinationsNotifierBuffered(bufferSize *int, buffererConfigs *[]destinations.BufferConfig, convertToJSONL *bool,
	dataChannel chan models.EntryData, shutdownChannel chan struct{}, reloadChannel chan bool, readyChannel chan bool) {

	var buffer []models.EntryData
	currentBufferSize := 0
	wg := sync.WaitGroup{}

selectloop:
	for {
		select {
		case <-shutdownChannel:
			resetBuffer(buffer, *buffererConfigs, *convertToJSONL, &wg)

			wg.Wait()
			readyChannel <- true

			break selectloop

		case <-reloadChannel:

			go resetBuffer(buffer, *buffererConfigs, *convertToJSONL, &wg)

			buffer = buffer[:0]
			currentBufferSize = 0

			readyChannel <- true

		case data := <-dataChannel:

			buffer = append(buffer, data)
			currentBufferSize += len(data.Body)

			if currentBufferSize >= *bufferSize {
				resetBuffer(buffer, *buffererConfigs, *convertToJSONL, &wg)
				buffer = buffer[:0]
				currentBufferSize = 0
			}
		}
	}
}

func resetBuffer(buffer []models.EntryData, buffererConfigs []destinations.BufferConfig, convertToJSONL bool, wg *sync.WaitGroup) {
	for _, destination := range buffererConfigs {
		wg.Add(1)
		go func(d destinations.BufferConfig) {
			defer wg.Done()
			if err := d.Notify(buffer, convertToJSONL); err != nil {
				destinationType := reflect.TypeOf(d).String()
				logging.Logger.Sugar().Errorf("Error notifying bufferer %s, \n%s", err, destinationType)
			}
		}(destination)
	}
}
