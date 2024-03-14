package destbuffer

import (
	"buffer-handler/logging"
	"buffer-handler/models"
	"buffer-handler/utils"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"
)

type FileBufferConfigUnique struct {
	UniqueFilePerBuffer bool
	FilePathFormat      string
	FileFormat          string
	ItemSeparator       string
	FileExtension       string
}

func (cfg FileBufferConfigUnique) Notify(entries []models.EntryData, convertToJSONL bool) error {
	logging.Logger.Debug("Start Buffering Unique")

	filePath, fileName = cfg.generateFilePathAndName()
	randomNumber := fmt.Sprintf("%03d", rand.Int31n(1000))
	fileName = fmt.Sprintf("%s%s.%s", fileName, randomNumber, cfg.FileExtension)

	err := utils.CreateFolderIfNotExist(filePath)
	if err != nil {
		logging.Logger.Sugar().Fatalf("Couldn't find and create path %s\n Error: %s", filePath, err)
	}

	filePath = path.Join(filePath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, entry := range entries {
		var data string
		if convertToJSONL {
			buffer := new(bytes.Buffer)
			if err := json.Compact(buffer, entry.Body); err != nil {
				logging.Logger.Error(err.Error())
				continue
			}
			data = buffer.String()
		} else {
			data = string(entry.Body)
		}

		_, err := writer.WriteString(data + cfg.ItemSeparator)
		if err != nil {
			return err
		}
	}

	logging.Logger.Sugar().Debugf("Data written to file: %s\n", filePath)
	return nil
}

func (cfg FileBufferConfigUnique) generateFilePathAndName() (string, string) {
	currentTime := time.Now()

	filePath := utils.FormatConfig(cfg.FilePathFormat, currentTime)
	fileName := utils.FormatConfig(cfg.FileFormat, currentTime)

	return filePath, fileName
}

func (cfg FileBufferConfigUnique) String() string {
	filePath, fileName := cfg.generateFilePathAndName()
	return path.Join(filePath, fmt.Sprintf("%s.%s", fileName, cfg.FileExtension))
}
