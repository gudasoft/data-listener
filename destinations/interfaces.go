package destinations

import "buffer-handler/models"

type StreamConfig interface {
	Notify(models.EntryData) error
}

type BufferConfig interface {
	Notify([]models.EntryData, bool) error
}
