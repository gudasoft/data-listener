package stream

import (
	"datalistener/src/logging"
	"datalistener/src/models"
	"datalistener/src/utils"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"
)

type FileStreamConfigUnique struct {
	UniqueFilePerRequest bool
	FilePathFormat       string
	FileFormat           string
	ItemSeparator        string
	FileExtension        string
}

func (cfg FileStreamConfigUnique) Notify(entry models.EntryData) error {
	logging.Logger.Debug("File Streaming Unique")

	filePath, fileName = cfg.generateFilePathAndName()
	randomNumber := fmt.Sprintf("%03d", rand.Int31n(1000))
	fileName = fmt.Sprintf("%s%s.%s", fileName, randomNumber, cfg.FileExtension)

	err := utils.CreateFolderIfNotExist(filePath)
	if err != nil {
		logging.Logger.Sugar().Fatalf("Couldn't find and create path %s\n Error: %s", filePath, err)
	}

	filePath = path.Join(filePath, fileName)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logging.Logger.Sugar().Fatalf("Problem with file opening: %s", err.Error())
	}
	defer file.Close()

	_, err = file.WriteString(string(entry.Body))
	if err != nil {
		return err
	}

	logging.Logger.Sugar().Debug("Data written to file: %s\n", filePath)
	return nil
}

func (cfg FileStreamConfigUnique) generateFilePathAndName() (string, string) {
	currentTime := time.Now()

	filePath := utils.FormatConfig(cfg.FilePathFormat, currentTime)
	fileName := utils.FormatConfig(cfg.FileFormat, currentTime)

	return filePath, fileName
}

func (cfg FileStreamConfigUnique) String() string {
	filePath, fileName := cfg.generateFilePathAndName()
	return path.Join(filePath, fmt.Sprintf("%s.%s, uniqueFilePerRequest: true", fileName, cfg.FileExtension))
}
