package buffer

import (
	"bufio"
	"bytes"
	"datalistener/src/logging"
	"datalistener/src/models"
	"datalistener/src/utils"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"
)

type FileBufferConfig struct {
	UniqueFilePerBuffer bool
	MaxFileSize         int
	FilePathFormat      string
	FileFormat          string
	ItemSeparator       string
	FileExtension       string
}

var fileSizeInMegabytes int
var filePath string
var fileName string

func (cfg FileBufferConfig) Notify(entries []models.EntryData, convertToJSONL bool) error {
	logging.Logger.Debug("Start Buffering")

	if fileSizeInMegabytes > cfg.MaxFileSize || fileSizeInMegabytes == 0 {
		filePath, fileName = cfg.generateFilePathAndName()
		randomNumber := fmt.Sprintf("%03d", rand.Int31n(1000))
		fileName = fmt.Sprintf("%s%s.%s", fileName, randomNumber, cfg.FileExtension)

		err := utils.CreateFolderIfNotExist(filePath)
		if err != nil {
			logging.Logger.Sugar().Fatalf("Couldn't find and create path %s\n Error: %s", filePath, err)
		}

		filePath = path.Join(filePath, fileName)
		fileSizeInMegabytes = 0
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, entry := range entries {
		var data string
		fileSizeInMegabytes += len(entry.Body)
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

	logging.Logger.Sugar().Debug("Data written to file: %s\n", filePath)
	return nil
}

func (cfg FileBufferConfig) generateFilePathAndName() (string, string) {
	currentTime := time.Now()

	filePath := utils.FormatConfig(cfg.FilePathFormat, currentTime)
	fileName := utils.FormatConfig(cfg.FileFormat, currentTime)

	return filePath, fileName
}

func (cfg FileBufferConfig) String() string {
	filePath, fileName := cfg.generateFilePathAndName()
	return path.Join(filePath, fmt.Sprintf("%s.%s, uniqueFilePerBuffer: %t,  maxFileSize: %dkB", fileName, cfg.FileExtension, cfg.UniqueFilePerBuffer, cfg.MaxFileSize/1024))
}
