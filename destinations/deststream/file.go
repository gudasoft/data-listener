package deststream

import (
	"buffer-handler/logging"
	"buffer-handler/models"
	"buffer-handler/utils"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"
)

type FileStreamConfig struct {
	UniqueFilePerRequest bool
	MaxFileSize          int
	FilePathFormat       string
	FileFormat           string
	ItemSeparator        string
	FileExtension        string
}

var fileSizeInMegabytes int
var filePath string
var fileName string

func (cfg FileStreamConfig) Notify(entry models.EntryData) error {
	logging.Logger.Debug("File Streaming")

	if cfg.UniqueFilePerRequest || fileSizeInMegabytes == 0 || fileSizeInMegabytes > cfg.MaxFileSize {
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
		logging.Logger.Sugar().Fatalf("Problem with file opening: %s", err.Error())
	}
	defer file.Close()

	fileSizeInMegabytes += len(entry.Body)
	_, err = file.WriteString(string(entry.Body) + cfg.ItemSeparator)
	if err != nil {
		return err
	}

	logging.Logger.Sugar().Debug("Data written to file: %s\n", filePath)

	return nil
}

func (cfg FileStreamConfig) generateFilePathAndName() (string, string) {
	currentTime := time.Now()

	filePath := utils.FormatConfig(cfg.FilePathFormat, currentTime)
	fileName := utils.FormatConfig(cfg.FileFormat, currentTime)

	return filePath, fileName
}

func (cfg FileStreamConfig) String() string {
	filePath, fileName := cfg.generateFilePathAndName()
	return path.Join(filePath, fmt.Sprintf("%s.%s, uniqueFilePerRequest: %t, maxFileSize: %dMB", fileName, cfg.FileExtension, cfg.UniqueFilePerRequest, cfg.MaxFileSize/1024/1024))
}
