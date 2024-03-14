package destbuffer

import (
	"buffer-handler/models"
	"fmt"
)

type DebugBuffer struct {
}

func (s DebugBuffer) Notify(e []models.EntryData) error {

	fmt.Println(e)
	return nil
}
