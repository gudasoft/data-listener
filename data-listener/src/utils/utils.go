package utils

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var LayoutMapping map[string]string

func init() {
	LayoutMapping = map[string]string{
		"%Y":  "2006",
		"%MO": "01",
		"%D":  "02",
		"%H":  "15",
		"%MN": "04",
		"%S":  "05",
		"%MS": "999",
	}
}

func FormatConfig(format string, currentTime time.Time) string {

	for key, value := range LayoutMapping {
		if key == "%MS" {
			format = strings.ReplaceAll(format, key, fmt.Sprintf("%03d", (currentTime.UnixNano()/int64(time.Millisecond))%1000))
		} else {
			format = strings.ReplaceAll(format, key, currentTime.Format(value))
		}
	}
	return format
}

func CreateFolderIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
