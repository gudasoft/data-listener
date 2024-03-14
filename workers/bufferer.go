package workers

import (
	"buffer-handler/destinations"
	"buffer-handler/models"
)

var (
	buffer            []models.EntryData
	currentBufferSize int
)

func DestinationsNotifierBuffered(bufferSize *int, buffererConfigs *[]destinations.BufferConfig, convertToJSONL *bool,
	ch chan models.EntryData, shutdownChan chan struct{}, reloadChan chan bool) {
	currentBufferSize = 0

	go func() {
		<-shutdownChan
		resetBuffer(buffer, *buffererConfigs, *convertToJSONL)
	}()

	go func() {
		for range reloadChan {
			resetBuffer(buffer, *buffererConfigs, *convertToJSONL)
			buffer = nil
			currentBufferSize = 0
		}
	}()

	for data := range ch {
		buffer = append(buffer, data)
		currentBufferSize += len(data.Body)

		if currentBufferSize >= *bufferSize {
			resetBuffer(buffer, *buffererConfigs, *convertToJSONL)
			buffer = nil
			currentBufferSize = 0
		}
	}

	resetBuffer(buffer, *buffererConfigs, *convertToJSONL)
}

func resetBuffer(buffer []models.EntryData, buffererConfigs []destinations.BufferConfig, convertToJSONL bool) {
	for _, destination := range buffererConfigs {
		destination.Notify(buffer, convertToJSONL)
	}
}

/*

func DestinationsNotifierBuffered(bufferSize *int, buffererConfigs *[]destinations.BufferConfig, convertToJSONL *bool,
	dataChannel chan models.EntryData, shutdownChannel chan struct{}, reloadChannel chan bool) {
	currentBufferSize = 0

selectloop:
	for {
		select {
		case <-shutdownChannel:
			resetBuffer(buffer, *buffererConfigs, *convertToJSONL)
			break selectloop

		case <-reloadChannel:
			go resetBuffer(buffer, *buffererConfigs, *convertToJSONL)
			buffer = buffer[:0]
			currentBufferSize = 0

		case data := <-dataChannel:

			buffer = append(buffer, data)
			currentBufferSize += len(data.Body)

			if currentBufferSize >= *bufferSize {
				resetBuffer(buffer, *buffererConfigs, *convertToJSONL)
				buffer = buffer[:0]
				currentBufferSize = 0
			}
		}
	}
}

*/

// :TODO Those two lines have have name: resetBuffer()
// there will be no penalty for calling that naming function because it will be inlined
// https://stackoverflow.com/questions/45836981/is-it-possible-to-inline-function-containing-loop-in-golang
//
// Using a for loop with a select statement is a common and efficient pattern
// in Go for concurrent programming, especially when you need to continuously
// monitor multiple channels for incoming data or events. It's an effective way
// to handle concurrency and manage goroutines efficiently. Here's why it's considered efficient:

// Non-blocking: The select statement doesn't block your program; instead, it actively
// listens to all the specified channels. When one of the channels becomes ready, it executes
// the corresponding case block. This means that your code won't waste CPU cycles by polling
// or busy-waiting on the channels.

// Concurrency: By using a for loop with a select statement, you can easily manage multiple
// goroutines and allow them to run concurrently. This is especially useful for scenarios
// where you want to handle multiple channels concurrently and respond to events or data as they
// become available.

// Scalability: The select statement can be used with a large number of channels, making it suitable
// for scalable solutions that involve managing many concurrent operations. You can efficiently add
// or remove channels from the select statement without introducing significant overhead.

// Fairness: The select statement randomly selects one of the ready cases. This fairness property
// ensures that if multiple channels become ready simultaneously, they have an equal chance
// of being selected. This can help prevent favoring one channel over others in scenarios with
// multiple concurrent inputs.

// While using a for loop with a select statement is efficient and idiomatic in Go for many concurrent
// programming tasks, it's important to design your program carefully to avoid potential issues
// like deadlock or livelock. Proper synchronization and communication between goroutines are key
// to building robust and efficient concurrent programs in Go.
